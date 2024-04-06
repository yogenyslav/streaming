import pymongo
import asyncio
import json
import base64
from aiokafka import AIOKafkaConsumer, AIOKafkaProducer
from PIL import Image
from io import BytesIO
from ultralytics import YOLO
from ultralytics.engine.results import Results


client = pymongo.MongoClient("mongodb://frame_mongo:27017/dev")
processing: dict[int, bool] = {}


def inference(model: YOLO, file: str, filename: str):
    img = Image.open(BytesIO(base64.b64decode(file)))
    results: list[Results] = model(img)
    for res in results:
        print(f"time {sum(res.speed.values())/1000}")
        if res.boxes.data.int().numpy().shape[0] == 0:
            continue
        data = res.boxes.data.int().numpy()
        print(data)
        for d in data:
            client.get_database("dev").get_collection("boxes").insert_one(
                {
                    "filename": filename,
                    "lb": int(d[0]),
                    "lt": int(d[1]),
                    "rb": int(d[2]),
                    "rt": int(d[3]),
                }
            )


async def detect():
    model = YOLO("yolov9c.pt")
    consumer = AIOKafkaConsumer(
        "frames",
        "cancel",
        bootstrap_servers="kafka:29092",
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
    )

    producer = AIOKafkaProducer(
        bootstrap_servers="kafka:29092",
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
        acks="all",
        enable_idempotence=True,
    )
    await producer.start()
    await consumer.start()
    try:
        async for msg in consumer:
            query_id: int = msg.value["query_id"]
            print(f"\nquery_id {query_id}")

            if query_id not in processing:
                processing[query_id] = True

            match msg.topic:
                case "cancel":
                    processing[query_id] = False
                case "frames":
                    filename: str = msg.value["filename"]
                    total_frames: int = msg.value["total_frames"]
                    file: str = msg.value["data"]
                    print(f"filename {filename}")
                    print(f"total_frames {total_frames}")

                    if query_id in processing and processing[query_id]:
                        inference(model, file, filename)
                        await producer.send_and_wait(
                            f"status_{query_id}",
                            {"filename": filename, "status": "success"},
                        )
                        print(f"frame {filename} processed")

                        if int(filename.split("_")[-1].split(".")[0]) == total_frames:
                            processing[query_id] = False
                            print("done")

    except Exception as e:
        print(str(e))
        await producer.send_and_wait(
            f"status_{query_id}", {"filename": filename, "status": "error"}
        )
        processing[query_id] = False

    finally:
        await producer.stop()
        await consumer.stop()


if __name__ == "__main__":
    asyncio.run(detect())
