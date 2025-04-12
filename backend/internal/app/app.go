package app

import (
	"backend/internal/controller/http/v1/website"
	bookTask "backend/internal/domain/usecases/book"
	commentTask "backend/internal/domain/usecases/comment"
	drawingTask "backend/internal/domain/usecases/drawing"
	bookRepo "backend/internal/repository/task/book"
	commentRepo "backend/internal/repository/task/comment"
	drawingRepo "backend/internal/repository/task/drawing"
	"backend/pkg/config"
	"backend/pkg/connectors/S3connector"
	"backend/pkg/connectors/pgconnector"
	"backend/pkg/connectors/rabbitconnector"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeOutShutdownService = time.Duration(5) * time.Second
	timeReadHeader         = time.Duration(5) * time.Second
)

func NewApp(config config.AppConfig) error {

	// Инициализация PostgreSQL
	connectorConfig, err := pgconnector.CreateConfig(&pgconnector.ConnectionConfig{
		Host:     config.PGStorage.Host,
		Port:     fmt.Sprint(config.PGStorage.Port),
		User:     config.PGStorage.User,
		Password: config.PGStorage.Pass,
		DbName:   config.PGStorage.DB,
		SslMode:  "disable",
	},
		nil)

	pgStorage, err := pgconnector.NewPgConnector(
		connectorConfig,
		10*time.Second,
		10*time.Second,
	)
	if err != nil {
		return fmt.Errorf("initialize pg storage: %w", err)
	}
	defer func() {
		pgStorage.CloseConnection()
	}()

	// Инициализация S3
	zap.L().Info("initialize S3 Storage...")
	s3Connector, err := S3connector.NewS3Client(S3connector.S3Config{
		Endpoint:  config.S3.Endpoint,
		Region:    config.S3.Region,
		AccessKey: config.S3.AccessKey,
		SecretKey: config.S3.SecretKey,
		Bucket:    config.S3.BucketName,
	})
	if err != nil {
		return fmt.Errorf("initialize s3 client: %w", err)
	}

	// Инициализация RabbitMQ
	zap.L().Info("Initializing rabbitMQ...")
	rabbitMQ, err := rabbitconnector.NewConnector(&rabbitconnector.RabbitConfig{
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		Username: config.RabbitMQ.Username,
		Password: config.RabbitMQ.Password,
		Path:     config.RabbitMQ.Path,
	})
	if err != nil {
		return fmt.Errorf("initialize rabbitMQ: %w", err)
	}
	defer func() {
		if err = rabbitMQ.CloseConnection(); err != nil {
			zap.L().Warn("failed to close rabbitMQ connection", zap.Error(err))
		}
	}()

	// Настройка роутера
	router := gin.New()

	if config.Debug {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.Use(otelgin.Middleware("task-service"))
	routersInit(
		router,
		drawingTask.NewUseCase(drawingRepo.NewRepository(pgStorage), rabbitMQ, s3Connector, config.S3.BucketName, config.Queues.InfoMessage),
		commentTask.NewUseCase(commentRepo.NewRepository(pgStorage), rabbitMQ, config.Queues.InfoMessage),
		bookTask.NewUseCase(bookRepo.NewRepository(pgStorage), rabbitMQ, config.Queues.InfoMessage),
		config.HTTPServer,
	)

	srv := &http.Server{
		Addr:        config.HTTPServer.Address,
		Handler:     router,
		ReadTimeout: timeReadHeader,
	}

	go func() {
		zap.L().Info("Server is starting", zap.String("address", "http://"+config.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), timeOutShutdownService)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("emergency shutdown http server", zap.Error(err))
	}
	zap.L().Info("http server shutdown")
	return nil
}

func routersInit(
	router *gin.Engine,
	drawingCase website.DrawingUsecase,
	commentCase website.CommentUsecase,
	bookCase website.BookUsecase,
	srv config.HttpServer,
) {
	website.Router(router.Group("cars_library"), drawingCase, commentCase, bookCase, srv.User, srv.Pass)
}
