package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"streaming/config"
	servresp "streaming/internal/server/response"
	"streaming/internal/video"
	"streaming/internal/video/query"
	"streaming/internal/video/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggermw "github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/logger"
	"github.com/yogenyslav/storage/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	cfg *config.Config
	app *fiber.App
	pg  *pgxpool.Pool
}

func New(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		BodyLimit:    1024 * 1024 * 1024,
		ErrorHandler: servresp.ErrorHandler,
		AppName:      "Streaming API",
	})
	app.Use(loggermw.New())
	app.Use(recovermw.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CorsOrigins,
		AllowHeaders:     "",
		AllowCredentials: false,
	}))

	pg := postgres.MustNew(&cfg.Postgres, 20)

	return &Server{
		cfg: cfg,
		app: app,
		pg:  pg,
	}
}

func (s *Server) Run() {
	s.app.Static("/static", "/static")

	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	frameAddr := fmt.Sprintf("%s:%d", s.cfg.FrameService.Host, s.cfg.FrameService.Port)
	frameConn, err := grpc.Dial(frameAddr, grpcOpts...)
	if err != nil {
		logger.Panicf("failed to connect to searchEngine: %v", err)
	}

	responseRepo := response.NewRepo(s.pg)
	responseController := response.NewController(responseRepo, frameConn)
	responseHandler := response.NewHandler(responseController)
	video.SetupResponseRoutes(s.app, responseHandler)

	queryRepo := query.NewRepo(s.pg)
	queryController := query.NewController(queryRepo, responseController, frameConn)
	queryHandler := query.NewHandler(queryController)
	video.SetupQueryRoutes(s.app, queryHandler)

	go s.listen(&s.cfg.Server)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	s.pg.Close()
	if err = s.app.Shutdown(); err != nil {
		logger.Warnf("failed to shutdown the app: %v", err)
	}
	if err = frameConn.Close(); err != nil {
		logger.Warnf("failed to close frameService grpc conn: %v", err)
	}
	logger.Info("app was gracefully shutdown")
	os.Exit(0)
}

func (s *Server) listen(cfg *config.ServerConfig) {
	if err := s.app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		logger.Errorf("unexpectedly stopping the server: %v", err)
	}
}
