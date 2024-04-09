import cv2
import minio.objectlockconfig
import pymongo
import gridfs
import asyncio
import pb.frame_pb2_grpc
import numpy as np
import grpc
import base64
import minio
from PIL import Image
from io import BytesIO
from detection.detect import scheduler, outbox
from bson.objectid import ObjectId
from pb.frame_pb2 import (
    Query,
    Response,
    ResponseStatus,
    ProcessedResp,
    ProcessedReq,
    QueryType,
)
from grpc import ServicerContext
from detection import detect
from concurrent import futures


client = pymongo.MongoClient("mongodb://mongo:27017/dev")
fs = gridfs.GridFS(client.get_database("dev"))
s3 = minio.Minio(
    "minio:9000",
    access_key="ABYVB5wEg4PDnIx0U7DU",
    secret_key="I56vUpDXJAgK8Tjlp9jPoTvelTzNnhKrtwEVLxmX",
    secure=False,
)

if not s3.bucket_exists("frame"):
    s3.make_bucket("frame", "eu-west-1", True)


class FrameService(pb.frame_pb2_grpc.FrameServiceServicer):
    async def Process(self, query: Query, context: grpc.aio.ServicerContext):
        print(f"query.source = {query.source} query.type = {query.type}")

        if query.type == QueryType.File:
            url = s3.presigned_get_object("video", query.source)
            print(f"url = {url}")
            cap = cv2.VideoCapture(url)
        else:
            cap = cv2.VideoCapture(f"{query.source}")

        frame = 0
        count = 0
        fps = int(cap.get(cv2.CAP_PROP_FPS))
        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        while True:
            if context.cancelled():
                print("context was canceled")
                await detect.cancel(query.id)
                scheduler.remove_all_jobs()
                return Response(status=ResponseStatus.Canceled)

            _, image = cap.read()

            img_str = base64.b64encode(cv2.imencode(".jpg", image)[1].tobytes())
            if frame % (fps // 2) == 0:
                count += 1
                filename = f"{query.id}_{count}.jpg"
                fs.put(
                    img_str,
                    encoding="utf-8",
                    filename=filename,
                    query_id=query.id,
                    sent=False,
                    processed=False,
                    total_frames=(total_frames // (fps // 2)) + 1,
                )

                print(f"frame {count}")

            frame += 1
            if frame >= total_frames:
                break

        status = await detect.receive(
            query.id, (total_frames // (fps // 2)) + 1, context
        )
        print(f"finished processing with status {status}")
        res: ResponseStatus
        match status:
            case "success":
                res = ResponseStatus.Success
            case "error":
                res = ResponseStatus.Error
            case "canceled":
                res = ResponseStatus.Canceled
        return Response(status=res)

    def FindProcessed(self, query: ProcessedReq, context: ServicerContext):
        files = (
            client.get_database("dev")
            .get_collection("fs.files")
            .find({"query_id": query.queryId, "processed": True})
        )

        src: list[str] = []

        for f in files:
            boxes = (
                client.get_database("dev")
                .get_collection("boxes")
                .find({"filename": f["filename"]})
            )

            flag = False
            res = fs.get(ObjectId(f["_id"]))
            path = f"processed_{f['filename']}"
            t = base64.b64decode(res.read().decode("utf-8"))
            img = np.array(Image.open(BytesIO(t)))
            for box in boxes:
                flag = True
                cv2.rectangle(
                    img, (box["lb"], box["lt"]), (box["rb"], box["rt"]), (255, 0, 0), 1
                )

            if flag:
                data = BytesIO()
                Image.fromarray(img).save(data, "JPEG")
                data.seek(0)
                s3.put_object("frame", path, data, len(data.getvalue()))
                res = s3.presigned_get_object("frame", path)
                src.append(res)
        return ProcessedResp(src=src)


async def serve():
    s = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))
    pb.frame_pb2_grpc.add_FrameServiceServicer_to_server(FrameService(), s)
    s.add_insecure_port("[::]:10000")

    scheduler.add_job(
        outbox,
        "interval",
        seconds=3,
        executor="default",
    )
    scheduler.start()
    await s.start()
    await s.wait_for_termination()


if __name__ == "__main__":
    asyncio.run(serve())
