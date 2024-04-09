import json
import pymongo
import gridfs
import base64
import asyncio
import cv2
import numpy as np
import minio
from io import BytesIO
from PIL import Image
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from bson import ObjectId
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from grpc.aio import ServicerContext
from apscheduler.executors.pool import ThreadPoolExecutor, ProcessPoolExecutor


_client = pymongo.MongoClient("mongodb://mongo:27017/dev")
executors = {"default": ThreadPoolExecutor(20), "processpool": ProcessPoolExecutor(5)}

scheduler = AsyncIOScheduler(executors=executors)
_fs = gridfs.GridFS(_client.get_database("dev"))

_s3 = minio.Minio(
    "minio:9000",
    access_key="ABYVB5wEg4PDnIx0U7DU",
    secret_key="I56vUpDXJAgK8Tjlp9jPoTvelTzNnhKrtwEVLxmX",
    secure=False,
)


def outbox():
    col = _client.get_database("dev").get_collection("fs.files")
    files = col.find({"sent": False})
    for file in files:
        data = _fs.get(ObjectId(file["_id"]))
        success = asyncio.run(
            process(
                {
                    "query_id": file["query_id"],
                    "filename": file["filename"],
                    "data": data.read().decode("utf-8"),
                    "total_frames": file["total_frames"],
                }
            )
        )
        print(f"{file['filename']} {success}")
        if success:
            col.update_one({"_id": file["_id"]}, {"$set": {"sent": True}})


async def process(data: dict) -> bool:
    success = False
    producer = AIOKafkaProducer(
        bootstrap_servers="kafka-service:29092",
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
        acks="all",
        enable_idempotence=True,
    )
    await producer.start()
    try:
        await producer.send_and_wait("frames", data)
        success = True
    except Exception as e:
        print(str(e))
    finally:
        await producer.stop()
        return success


async def cancel(query_id: int):
    producer = AIOKafkaProducer(
        bootstrap_servers="kafka-service:29092",
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
        acks="all",
        enable_idempotence=True,
    )
    await producer.start()
    try:
        await producer.send_and_wait("cancel", {"query_id": query_id})
    finally:
        await producer.stop()


async def receive(query_id: int, total_frames: int, context: ServicerContext):
    consumer = AIOKafkaConsumer(
        f"status_{query_id}",
        bootstrap_servers="kafka-service:29092",
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
    )

    col = _client.get_database("dev").get_collection("fs.files")
    processed = 0

    await consumer.start()
    try:
        async for msg in consumer:
            processed += 1
            if context.cancelled():
                print("context was canceled")
                await cancel(query_id)
                await consumer.stop()
                return "canceled"
            else:
                status = msg.value["status"]
                filename = msg.value["filename"]
                print(f"received {filename} {status}")

                if status == "success":
                    f = (
                        _client.get_database("dev")
                        .get_collection("fs.files")
                        .find_one({"filename": filename})
                    )

                    boxes = (
                        _client.get_database("dev")
                        .get_collection("boxes")
                        .find({"filename": filename})
                    )

                    flag = False
                    res = _fs.get(ObjectId(f["_id"]))
                    path = f"processed_{filename}"
                    t = base64.b64decode(res.read().decode("utf-8"))
                    img = np.array(Image.open(BytesIO(t)))
                    for box in boxes:
                        flag = True
                        cv2.rectangle(
                            img,
                            (box["lb"], box["lt"]),
                            (box["rb"], box["rt"]),
                            (255, 0, 0),
                            1,
                        )

                    if flag:
                        col.update_one(
                            {"filename": filename},
                            {"$set": {"processed": True}},
                        )
                        data = BytesIO()
                        Image.fromarray(img).save(data, "JPEG")
                        data.seek(0)
                        _s3.put_object("frame", path, data, len(data.getvalue()))
                        res = _s3.presigned_get_object("frame", path)

                if processed == total_frames:
                    return status
    finally:
        await consumer.stop()
