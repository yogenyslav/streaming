package response

import (
	"context"
	"io"
	"net/http"
	"time"

	"streaming/internal/pb"
	"streaming/internal/shared"
	"streaming/internal/video/response/model"
	"streaming/pkg/storage/minios3"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/logger"
	"google.golang.org/grpc"
)

type responseRepo interface {
	InsertOne(ctx context.Context, params model.Response) error
	UpdateOne(ctx context.Context, params model.Response) error
	FindOneByQueryId(ctx context.Context, id int64) (model.Response, error)
}

type Controller struct {
	repo         responseRepo
	frameService pb.FrameServiceClient
	s3           *minios3.S3
}

func NewController(repo responseRepo, frameConn *grpc.ClientConn, s3 *minios3.S3) *Controller {
	return &Controller{
		repo:         repo,
		frameService: pb.NewFrameServiceClient(frameConn),
		s3:           s3,
	}
}

func (ctrl *Controller) InsertOne(ctx context.Context, params model.ResponseCreateReq) error {
	resp := model.Response{
		QueryId: params.QueryId,
	}
	if err := ctrl.repo.InsertOne(ctx, resp); err != nil {
		logger.Errorf("failed to insert response: %v", err)
		return shared.ErrInsertRecord
	}

	return nil
}

func (ctrl *Controller) UpdateOne(ctx context.Context, params model.ResponseUpdateReq) error {
	resp := model.Response{
		QueryId: params.QueryId,
		Status:  params.Status,
	}
	if err := ctrl.repo.UpdateOne(ctx, resp); err != nil {
		logger.Errorf("failed to update response: %v", err)
		return shared.ErrUpdateRecord
	}
	return nil
}

func (ctrl *Controller) FindOneByQueryId(ctx context.Context, queryId int64) (model.ResponseDto, error) {
	var res model.ResponseDto

	response, err := ctrl.repo.FindOneByQueryId(ctx, queryId)
	if err != nil {
		logger.Errorf("failed to find response: %v", err)
		return res, shared.ErrFindRecord
	}

	in := &pb.ProcessedReq{QueryId: queryId}
	resp, err := ctrl.frameService.FindProcessed(ctx, in)
	if err != nil {
		logger.Errorf("failed to find processed frames: %v", err)
		return res, shared.ErrFindProcessed
	}

	res = response.ToDto()
	res.Source = resp.GetSrc()

	return res, nil
}

func (ctrl *Controller) GetStatic(ctx context.Context, name string) ([]byte, error) {
	url, err := ctrl.s3.PresignedGetObject(ctx, "frame", name, time.Hour*24*7, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("status is %s", resp.Status)
		return nil, fiber.ErrBadRequest
	}

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}
