package main

import (
	"backend/internal/app"
	"backend/pkg/config"
	"backend/pkg/utils"
	"log"

	"go.uber.org/zap"
)

func main() {

	var cfg config.AppConfig
	logger, err := utils.CreateLogger(zap.InfoLevel)
	if err != nil {
		log.Fatalf("creating logger failed: %v", err)
	}
	if err = cfg.ReadEnvConfig(); err != nil {
		logger.Fatal("reading environment variables failed", zap.Error(err))
	}
	if cfg.PathConfig != "" {
		if err = cfg.ReadYamlConfig(cfg.PathConfig); err != nil {
			logger.Fatal("reading config failed", zap.Error(err))
		}
	}

	if err = cfg.Validate(); err != nil {
		logger.Fatal("validating config failed", zap.Error(err))
	}
	if cfg.Debug {
		logger.Warn("application is running in debug mode")
		logger, err = utils.CreateLogger(zap.DebugLevel)
		if err != nil {
			log.Fatalf("failed to create logger: %s", err)
		}
	}
	zap.ReplaceGlobals(logger)

	if err = app.NewApp(cfg); err != nil {
		logger.Fatal("application failed", zap.Error(err))
		return
	}
}
