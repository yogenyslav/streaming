import json
import asyncio
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from grpc.aio import ServicerContext


async def process(data: dict):
    producer = AIOKafkaProducer(
        bootstrap_servers="kafka:29092",
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
    )
    await producer.start()
    try:
        await producer.send_and_wait("frames", data)
    finally:
        await producer.stop()


async def cancel(query_id: int):
    producer = AIOKafkaProducer(
        bootstrap_servers="kafka:29092",
        value_serializer=lambda x: json.dumps(x).encode(encoding="utf-8"),
    )
    await producer.start()
    try:
        await producer.send_and_wait("cancel", {"query_id": query_id})
    finally:
        await producer.stop()


async def receive(query_id: int, context: ServicerContext):
    consumer = AIOKafkaConsumer(
        "status",
        bootstrap_servers="kafka:29092",
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        max_poll_records=1,
    )

    status = ""

    await consumer.start()
    try:
        async for msg in consumer:
            if context.cancelled():
                print("context was canceled")
                await cancel(query_id)
                await consumer.stop()
                return "canceled"
            status = msg.value["status"]
            print(status)
            await consumer.stop()
            return status
    finally:
        await consumer.stop()
        return status
