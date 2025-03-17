package app

import (
	"backend/internal/controller/http/v1/website"
	"backend/internal/domain/usecases/task"
	task2 "backend/internal/repository/task/drawing"
	"backend/pkg/config"
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
		task.NewUseCase(task2.NewRepository(pgStorage), rabbitMQ, config.Queues.InfoMessage),
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
	usecase website.Usecase,
	srv config.HttpServer,
) {
	website.Router(router.Group("cars_library"), usecase, srv.User, srv.Pass)
}
