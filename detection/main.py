import os
import asyncio
import json
import base64
import aiokafka
from PIL import Image
from dotenv import load_dotenv
from io import BytesIO
from ultralytics import YOLO
from ultralytics.engine.results import Results

load_dotenv(".env")
processing: dict[int, bool] = {}


async def send(topic: str, data: dict):
    producer = aiokafka.AIOKafkaProducer(
        bootstrap_servers=os.getenv("KAFKA_HOST"),
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
        acks="all",
        enable_idempotence=True,
    )
    await producer.start()
    try:
        await producer.send_and_wait(topic, data)
    finally:
        await producer.stop()


async def inference(model: YOLO, file: str, query_id: int, frame_id: int, last: bool):
    img = Image.open(BytesIO(base64.b64decode(file)))
    results: list[Results] = model(img)
    for res in results:
        print(f"time {sum(res.speed.values())/1000}")
        if res.boxes.data.int().numpy().shape[0] == 0:
            if last:
                await send(
                    "responser",
                    {
                        "img": file,
                        "query_id": query_id,
                        "filename": f"{query_id}_{frame_id}.jpg",
                        "lb": -1,
                        "lt": -1,
                        "rb": -1,
                        "rt": -1,
                        "last": True,
                    },
                )
            continue
        data = res.boxes.data.int().numpy()
        print(data)
        for d in data:
            await send(
                "responser",
                {
                    "img": file,
                    "query_id": query_id,
                    "filename": f"{query_id}_{frame_id}.jpg",
                    "lb": int(d[0]),
                    "lt": int(d[1]),
                    "rb": int(d[2]),
                    "rt": int(d[3]),
                    "last": last,
                },
            )


async def detect():
    model = YOLO("yolov9c.pt")
    consumer = aiokafka.AIOKafkaConsumer(
        "detection",
        "detection-cancel",
        bootstrap_servers=os.getenv("KAFKA_HOST"),
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
    )

    await consumer.start()
    try:
        async for msg in consumer:
            try:
                print(f"Received message: {msg.offset}")
                cancel: bool = msg.value["cancel"]
                query_id: int = msg.value["query_id"]

                if msg.topic == "detection-cancel":
                    print(f"cancel {cancel}")
                    if cancel:
                        print(f"query_id {query_id} was canceled")
                        processing[query_id] = False
                        continue

                if query_id not in processing.keys():
                    processing[query_id] = True

                if not processing[query_id]:
                    continue

                frame_id: int = msg.value["frame_id"]
                img: str = msg.value["img"]
                last: bool = msg.value["last"]

                await inference(model, img, query_id, frame_id, last)

                print(f"Processed query {query_id} frame {frame_id}")

                if last:
                    print(f"Processed query {query_id} LAST frame {frame_id}")
                    await send(
                        "detection-res",
                        data={
                            "id": query_id,
                            "message": "sucess",
                        },
                    )
            except Exception as e:
                print(f"detection error {str(e)}")
                await send(
                    "detection-res",
                    {
                        "id": query_id,
                        "message": "error",
                    },
                )

    finally:
        await consumer.stop()


if __name__ == "__main__":
    asyncio.run(detect())
