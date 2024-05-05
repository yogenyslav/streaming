import aiokafka
import asyncio
import os
import pymongo
import cv2
import numpy as np
import json
import tempfile
import base64
from tempfile import TemporaryDirectory
from minio import Minio
from dotenv import load_dotenv
from PIL import Image
from io import BytesIO


load_dotenv(".env")


mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}/{os.getenv('MONGO_DB')}"
mongo_client = pymongo.MongoClient(mongo_url)
col = mongo_client.get_database(os.getenv("MONGO_DB")).get_collection("frames")

s3 = Minio(
    os.getenv("S3_HOST"),
    access_key=os.getenv("ACCESS_KEY"),
    secret_key=os.getenv("SECRET_KEY"),
    secure=False,
)

if not s3.bucket_exists("frames"):
    s3.make_bucket("frames", "eu-west-1", True)


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


async def main():
    consumer = aiokafka.AIOKafkaConsumer(
        "responser",
        bootstrap_servers=os.getenv("KAFKA_HOST"),
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
    )

    await consumer.start()
    try:
        async for msg in consumer:
            try:
                print(f"Received message: {msg.offset}")
                query_id = msg.value["query_id"]
                img: str = msg.value["img"]
                filename: str = msg.value["filename"]
                lb: int = msg.value["lb"]
                lt: int = msg.value["lt"]
                rb: int = msg.value["rb"]
                rt: int = msg.value["rt"]
                last: bool = msg.value["last"]

                with TemporaryDirectory() as temp_dir:
                    img = Image.open(BytesIO(base64.b64decode(img)))
                    img = np.array(img)
                    img = cv2.cvtColor(img, cv2.COLOR_RGB2BGR)
                    if lb != -1:
                        cv2.rectangle(img, (lb, lt), (rb, rt), (0, 255, 0), 2)
                    cv2.imwrite(f"{temp_dir}/{filename}", img)
                    s3.fput_object(
                        "frames",
                        filename,
                        f"{temp_dir}/{filename}",
                        content_type="image/jpeg",
                    )

                url = s3.presigned_get_object("frames", filename)
                col.insert_one(
                    {
                        "query_id": query_id,
                        "filename": filename,
                        "lb": lb,
                        "lt": lt,
                        "rb": rb,
                        "rt": rt,
                        "link": url,
                    }
                )

                print(f"Saved {filename} to s3 and mongo")

                if last:
                    await send(
                        "responser-res",
                        {
                            "id": query_id,
                            "message": "success",
                        },
                    )

            except Exception as e:
                print(f"responser error {str(e)}")

    finally:
        await consumer.stop()


if __name__ == "__main__":
    asyncio.run(main())
