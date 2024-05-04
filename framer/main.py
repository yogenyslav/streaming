import aiokafka
import asyncio
import os
import cv2
import json
import base64
from minio import Minio
from dotenv import load_dotenv


load_dotenv(".env")


s3 = Minio(
    os.getenv("S3_HOST"),
    access_key=os.getenv("ACCESS_KEY"),
    secret_key=os.getenv("SECRET_KEY"),
    secure=False,
)

processing: dict[int, bool] = {}


async def send_to_detection(data: dict, cancel: bool = False):
    data["cancel"] = cancel
    producer = aiokafka.AIOKafkaProducer(
        bootstrap_servers=os.getenv("KAFKA_HOST"),
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
        acks="all",
        enable_idempotence=True,
    )
    await producer.start()
    try:
        await producer.send_and_wait("detection", data)
    finally:
        await producer.stop()


async def send_framer_res(data: dict):
    producer = aiokafka.AIOKafkaProducer(
        bootstrap_servers=os.getenv("KAFKA_HOST"),
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
        acks="all",
        enable_idempotence=True,
    )
    await producer.start()
    try:
        await producer.send_and_wait("framer-res", data)
    finally:
        await producer.stop()


async def get_frames(id_: int, cap: cv2.VideoCapture):
    frame = 0
    cnt = 0
    fps = int(cap.get(cv2.CAP_PROP_FPS))
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    while True:
        if id_ not in processing.keys() or not processing[id_]:
            return

        frame += 1
        if frame == total_frames:
            cnt += 1
            _, image = cap.read()
            img_bytes = base64.b64encode(cv2.imencode(".jpg", image)[1].tobytes())
            await send_to_detection(
                data={
                    "query_id": id_,
                    "frame_id": cnt,
                    "img": img_bytes.decode(),
                    "last": True,
                }
            )

            await send_framer_res(
                data={
                    "id": id_,
                    "message": "success",
                }
            )

            return

        ok, image = cap.read()
        if ok and frame % (fps // 2) == 0:
            print(f"query {id_} frame {cnt}")
            cnt += 1
            img_bytes = base64.b64encode(cv2.imencode(".jpg", image)[1].tobytes())
            await send_to_detection(
                data={
                    "query_id": id_,
                    "frame_id": cnt,
                    "img": img_bytes.decode(),
                    "last": False,
                }
            )


async def main():
    consumer = aiokafka.AIOKafkaConsumer(
        "framer",
        bootstrap_servers=os.getenv("KAFKA_HOST"),
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
    )

    await consumer.start()
    try:
        async for msg in consumer:
            try:
                print(f"Received message: {msg.offset} {msg.key}")
                id_ = msg.value["id"]
                source = msg.value["source"]
                type_ = msg.value["type"]
                cancel = msg.value["cancel"]

                if cancel:
                    processing[id_] = False
                    await send_to_detection(
                        data={
                            "query_id": id_,
                            "frame_id": 0,
                            "img": "",
                            "cancel": True,
                        }
                    )
                    continue

                if id_ not in processing.keys():
                    processing[id_] = True

                if not processing[id_]:
                    continue

                if type_ == 1:
                    url = s3.presigned_get_object("video", source)
                    print(f"url = {url}")
                    cap = cv2.VideoCapture(url)
                else:
                    cap = cv2.VideoCapture(source)

                await get_frames(id_, cap)
            except Exception as e:
                print(f"framer error {str(e)}")
                await send_framer_res(
                    data={
                        "id": id_,
                        "message": "error",
                    }
                )

    finally:
        await consumer.stop()


if __name__ == "__main__":
    asyncio.run(main())
