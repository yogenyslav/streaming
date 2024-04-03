import logging
import cv2
import pymongo
import gridfs
import asyncio
import pb.frame_pb2_grpc
import base64
import numpy as np
import grpc
from bson.objectid import ObjectId
from pb.frame_pb2 import Query, Response, ResponseStatus, ProcessedResp, ProcessedReq
from grpc import ServicerContext, server
from detection import detect
from concurrent import futures


client = pymongo.MongoClient("mongodb://frame_mongo:27017/dev")
fs = gridfs.GridFS(client.get_database("dev"))


class FrameService(pb.frame_pb2_grpc.FrameServiceServicer):
    async def Process(self, query: Query, context: grpc.aio.ServicerContext):
        log.info(f"query = {query}")

        cap = cv2.VideoCapture(f".{query.source}")

        frame = 0
        count = 0
        fps = int(cap.get(cv2.CAP_PROP_FPS))
        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        while True:
            if context.cancelled():
                log.info("context was canceled")
                await detect.cancel(query.id)
                return Response(status=ResponseStatus.Canceled)

            _, image = cap.read()

            img_str = cv2.imencode(".jpg", image)[1].tostring()
            if frame % fps == 0:
                count += 1
                filename = f"{query.id}_{count}.jpg"
                fs.put(
                    img_str,
                    encoding="utf-8",
                    filename=filename,
                    query_id=query.id,
                )

                await detect.process(
                    {
                        "query_id": query.id,
                        "filename": filename,
                        "data": str(
                            base64.b64encode(cv2.imencode(".jpg", image)[1].tobytes()),
                        ),
                        "total_frames": total_frames // fps,
                    }
                )

                log.info(f"frame {count}")

            frame += 1
            if frame >= total_frames:
                break

        status = await detect.receive(query.id, context)
        log.info(f"finished processing with status {status}")
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
            .find({"query_id": query.queryId})
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
            path = f"/static/processed_{f['filename']}"
            img = cv2.imdecode(np.frombuffer(res.read(), np.uint8), cv2.IMREAD_COLOR)
            log.info(path)
            for box in boxes:
                flag = True
                cv2.rectangle(
                    img, (box["lb"], box["lt"]), (box["rb"], box["rt"]), (255, 0, 0), 1
                )

            if flag:
                cv2.imwrite(path, img)
                src.append(path)

        return ProcessedResp(src=src)


async def serve():
    s = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))
    pb.frame_pb2_grpc.add_FrameServiceServicer_to_server(FrameService(), s)
    s.add_insecure_port("[::]:10000")
    await s.start()
    await s.wait_for_termination()


logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

if __name__ == "__main__":
    asyncio.run(serve())
