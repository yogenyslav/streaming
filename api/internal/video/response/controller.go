package response

import (
	"context"

	"streaming/internal/pb"
	"streaming/internal/shared"
	"streaming/internal/video/response/model"

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
}

func NewController(repo responseRepo, frameConn *grpc.ClientConn) *Controller {
	return &Controller{
		repo:         repo,
		frameService: pb.NewFrameServiceClient(frameConn),
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
		return res, shared.ErrFindProcessed
	}

	res = response.ToDto()
	res.Source = resp.GetSrc()

	return res, nil
}
