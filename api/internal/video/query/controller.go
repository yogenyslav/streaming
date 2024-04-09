package query

import (
	"context"
	"mime/multipart"
	"strings"
	"sync"

	"streaming/internal/pb"
	"streaming/internal/shared"
	"streaming/internal/video/query/model"
	respModel "streaming/internal/video/response/model"
	"streaming/pkg"
	"streaming/pkg/storage/minios3"

	"github.com/minio/minio-go/v7"
	"github.com/yogenyslav/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queryRepo interface {
	InsertOne(ctx context.Context, params model.Query) (int64, error)
	UpdateSource(ctx context.Context, id int64, source string) error
}

type responseController interface {
	InsertOne(ctx context.Context, params respModel.ResponseCreateReq) error
	UpdateOne(ctx context.Context, params respModel.ResponseUpdateReq) error
}

type Controller struct {
	qr           queryRepo
	rc           responseController
	frameService pb.FrameServiceClient
	s3           *minios3.S3
	mu           sync.Mutex
}

func NewController(qr queryRepo, rc responseController, frameConn *grpc.ClientConn, s3 *minios3.S3) *Controller {
	return &Controller{
		qr:           qr,
		rc:           rc,
		s3:           s3,
		frameService: pb.NewFrameServiceClient(frameConn),
	}
}

func (ctrl *Controller) InsertOne(ctx context.Context, processing map[int64]context.CancelFunc, params model.QueryCreateReq) (int64, error) {
	query := model.Query{Type: pkg.QueryType(params.Source)}

	if query.Type == shared.TypeLink && params.Timeout == 0 {
		return 0, shared.ErrUnspecifiedTimeoutForLink
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

	go func() {
		var (
			rawFile     multipart.File
			err         error
			ctx, cancel = context.WithCancel(ctx)
		)
		defer cancel()

		ctrl.mu.Lock()
		processing[id] = cancel
		logger.Debugf("processing %v", processing)
		ctrl.mu.Unlock()

		defer func() {
			if err != nil {
				respUpdate := respModel.ResponseUpdateReq{
					QueryId: id,
					Status:  shared.StatusError,
				}
				if err = ctrl.rc.UpdateOne(context.Background(), respUpdate); err != nil {
					logger.Warnf("response status was not correctly updated: %v", err)
				}
			}
		}()

		if query.Type == shared.TypeFile {
			rawFile, err = params.File.Open()
			if err != nil {
				logger.Errorf("failed to open file: %v", err)
				return
			}

			split := strings.Split(params.File.Filename, ".")
			source := params.Name.String() + "." + split[len(split)-1]
			_, err = ctrl.s3.PutObject(ctx, shared.VideoBucket, source, rawFile, params.File.Size, minio.PutObjectOptions{})
			if err != nil {
				logger.Errorf("failed to put file into s3: %v", err)
				return
			}

			if err = ctrl.qr.UpdateSource(ctx, id, source); err != nil {
				logger.Errorf("failed to update file query source: %v", err)
				return
			}

			query.Source = source
		}

		in := &pb.Query{
			Id:     id,
			Type:   pb.QueryType(query.Type),
			Source: query.Source,
		}
		if params.Timeout == 0 {
			in.Timeout = nil
		} else {
			in.Timeout = &params.Timeout
		}

		if err = ctrl.process(ctx, processing, in); err != nil {
			return
		}
	}()

	return id, nil
}

func (ctrl *Controller) process(ctx context.Context, processing map[int64]context.CancelFunc, params *pb.Query) error {
	resp, err := ctrl.frameService.Process(ctx, params)
	if err != nil {
		grpcErr := status.Convert(err)
		if grpcErr.Code() != codes.Canceled {
			logger.Errorf("processing query %d failed: %v", params.Id, err)
			return shared.ErrProcessQuery
		} else {
			resp = &pb.Response{
				Status: pb.ResponseStatus_Canceled,
			}
		}
	}

	respUpdate := respModel.ResponseUpdateReq{
		QueryId: params.Id,
		Status:  shared.ResponseStatus(resp.GetStatus()),
	}
	if err = ctrl.rc.UpdateOne(context.Background(), respUpdate); err != nil {
		logger.Errorf("failed to update %d response status: %v", params.Id, err)
		return err
	}

	ctrl.mu.Lock()
	delete(processing, params.Id)
	logger.Debugf("processed %v", processing)
	ctrl.mu.Unlock()
	return nil
}

func (ctrl *Controller) CancelById(ctx context.Context, processing map[int64]context.CancelFunc, id int64) error {
	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	logger.Debugf("canceling %v", processing)
	cancel, ok := processing[id]
	if !ok {
		return shared.ErrQueryIsNotProcessed
	}
	cancel()

	delete(processing, id)
	logger.Debugf("canceled %v", processing)
	logger.Infof("processing query %d was canceled", id)
	return nil
}
