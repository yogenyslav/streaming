package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"streaming/api/config"
	srvconf "streaming/api/internal/server/config"
	"streaming/api/internal/streaming/pb"
	"streaming/api/internal/streaming/query"
	qc "streaming/api/internal/streaming/query/controller"
	qh "streaming/api/internal/streaming/query/handler"
	qr "streaming/api/internal/streaming/query/repo"
	"streaming/api/internal/streaming/response"
	rc "streaming/api/internal/streaming/response/controller"
	rh "streaming/api/internal/streaming/response/handler"
	rr "streaming/api/internal/streaming/response/repo"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/infrastructure/prom"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	srvresp "github.com/yogenyslav/pkg/response"
	"github.com/yogenyslav/pkg/storage/minios3"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	cfg    *config.Config
	db     postgres.Postgres
	s3     minios3.S3
	app    *fiber.App
	tracer trace.Tracer
}

func New(cfg *config.Config) *Server {
	errorHandler := srvresp.NewErrorHandler(errStatus)
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.Handler,
		BodyLimit:    cfg.Server.BodyLimit,
	})

	app.Use(logger.New())
	app.Use(recovermw.New())
	app.Use(otelfiber.Middleware())

	tracing.MustSetupOTel(fmt.Sprintf("%s:%d", cfg.Tracing.Host, cfg.Tracing.Port), "api")
	tracer := otel.Tracer("api")

	s3 := minios3.MustNew(cfg.S3, tracer)
	if err := s3.CreateBuckets(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed to create s3 buckets")
	}

	return &Server{
		cfg:    cfg,
		db:     postgres.MustNew(cfg.Postgres, tracer),
		s3:     s3,
		app:    app,
		tracer: tracer,
	}
}

func (s *Server) Run() {
	defer s.db.GetPool().Close()

	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	orchestratorAddr := fmt.Sprintf("%s:%d", s.cfg.Orchestrator.Host, s.cfg.Orchestrator.Port)
	conn, err := grpc.Dial(orchestratorAddr, grpcOpts...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to orchestrator")
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to properly close grpc connection")
		}
	}()

	responseRepo := rr.New(s.db)
	responseController := rc.New(responseRepo, pb.NewOrchestratorClient(conn), s.tracer)
	responseHandler := rh.New(responseController)
	response.SetupResponseRoutes(s.app, responseHandler)

	queryRepo := qr.New(s.db)
	queryController := qc.New(queryRepo, responseRepo, pb.NewOrchestratorClient(conn), s.s3, s.tracer)
	queryHandler := qh.New(queryController)
	query.SetupQueryRoutes(s.app, queryHandler)

	go s.listen(s.cfg.Server)
	go prom.HandlePrometheus(s.cfg.Prometheus.Host, s.cfg.Prometheus.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	if err = s.app.Shutdown(); err != nil {
		log.Warn().Err(err).Msg("failed to properly shutdown fiber.App")
	}

	log.Info().Msg("server was gracefully stopped")
	os.Exit(0)
}

func (s *Server) listen(cfg *srvconf.ServerConfig) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := s.app.Listen(addr); err != nil {
		log.Error().Err(err).Msg("error while serving http")
	}
}
