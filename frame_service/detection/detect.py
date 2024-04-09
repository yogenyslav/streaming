import json
import pymongo
import gridfs
import base64
import asyncio
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from bson import ObjectId
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from grpc.aio import ServicerContext
from apscheduler.executors.pool import ThreadPoolExecutor, ProcessPoolExecutor


_client = pymongo.MongoClient("mongodb://mongo:27017/dev")
executors = {"default": ThreadPoolExecutor(20), "processpool": ProcessPoolExecutor(5)}

scheduler = AsyncIOScheduler(executors=executors)
_fs = gridfs.GridFS(_client.get_database("dev"))


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
                print(f"received {msg.value['filename']} {status}")

                if status == "success":
                    col.update_one(
                        {"filename": msg.value["filename"]},
                        {"$set": {"processed": True}},
                    )
                if processed == total_frames:
                    return status
    finally:
        await consumer.stop()
