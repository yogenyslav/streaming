package query

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"streaming/api/internal/detection/query/model"
	respModel "streaming/api/internal/detection/response/model"
	"streaming/api/internal/pb"
	"streaming/api/internal/shared"
	"streaming/api/pkg"
	"streaming/api/pkg/storage/minios3"

	"github.com/minio/minio-go/v7"
	"github.com/yogenyslav/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queryRepo interface {
	InsertOne(ctx context.Context, params model.Query) (int64, error)
}

type responseController interface {
	InsertOne(ctx context.Context, params respModel.ResponseCreateReq) error
	UpdateOne(ctx context.Context, params respModel.ResponseUpdateReq) error
	FindOneByQueryId(ctx context.Context, queryId int64) (respModel.Response, error)
}

type Controller struct {
	qr           queryRepo
	rc           responseController
	orchestrator pb.OrchestratorClient
	s3           *minios3.S3
	mu           sync.Mutex
}

func NewController(qr queryRepo, rc responseController, orchestrator *grpc.ClientConn, s3 *minios3.S3) *Controller {
	return &Controller{
		qr:           qr,
		rc:           rc,
		s3:           s3,
		orchestrator: pb.NewOrchestratorClient(orchestrator),
	}
}

func (ctrl *Controller) InsertOne(ctx context.Context, params model.QueryCreateReq) (int64, error) {
	query := model.Query{
		Type:   params.Type,
		Source: params.Source,
	}

	if params.Type == shared.QueryTypeVideo {
		split := strings.Split(params.Name, ".")
		source := fmt.Sprintf("%s-%d.%s", params.Name, pkg.GetLocalTime().Unix(), split[len(split)-1])
		_, err := ctrl.s3.PutObject(ctx, shared.VideoBucket, source, params.Video, params.Size, minio.PutObjectOptions{})
		if err != nil {
			logger.Errorf("failed to put file into s3: %v", err)
			return 0, shared.ErrInsertRecord
		}
		query.Source = source
	}

	id, err := ctrl.qr.InsertOne(ctx, query)
	if err != nil {
		logger.Errorf("failed to insert query: %v", err)
		return 0, shared.ErrInsertRecord
	}

	respCreate := respModel.ResponseCreateReq{
		QueryId: id,
	}
	if err = ctrl.rc.InsertOne(ctx, respCreate); err != nil {
		return 0, err
	}

	go ctrl.process(context.Background(), id, query)

	return id, nil
}

func (ctrl *Controller) process(ctx context.Context, id int64, params model.Query) {
	in := &pb.Query{
		Id:     id,
		Source: params.Source,
		Type:   pb.QueryType(params.Type),
	}

	resp, err := ctrl.orchestrator.Process(ctx, in)
	if err != nil {
		grpcErr := status.Convert(err)
		if grpcErr.Code() != codes.Canceled {
			logger.Errorf("processing query %d failed: %v", id, err)
		}
	}

	logger.Infof("processing query %d finished with status %v", id, resp.GetStatus())

	if err = ctrl.rc.UpdateOne(ctx, respModel.ResponseUpdateReq{
		QueryId: id,
		Status:  shared.ResponseStatus(resp.GetStatus()),
	}); err != nil {
		logger.Errorf("failed to update response: %v", err)
	}
}

func (ctrl *Controller) CancelById(ctx context.Context, id int64) error {
	query, err := ctrl.rc.FindOneByQueryId(ctx, id)
	if err != nil {
		return err
	}

	if query.Status != shared.ResponseStatusCanceled {
		return shared.ErrQueryIsNotProcessed
	}

	resp, err := ctrl.orchestrator.Cancel(ctx, &pb.CancelReq{
		QueryId: id,
	})
	if err != nil || resp.GetSuccess() != true {
		logger.Errorf("failed to cancel query %d: %v", id, err)
		return err
	}
	return nil
}
