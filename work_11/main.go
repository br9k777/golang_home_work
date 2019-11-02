package main

import (
	"golang_home_work/work_11/config"
	"golang_home_work/work_11/daemon"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	var err error
	var cfg *config.Config
	config.InitDefaultConfigSettings()
	if cfg, err = config.GetConfig(); err != nil {
		zap.S().Fatal(err)
	}
	if err = viper.WriteConfigAs(".config.yaml"); err != nil {
		zap.S().Error(err)
	}

	var newLogger *zap.Logger
	if cfg.LogConfigurationFile != "" {
		if newLogger, err = config.GetLoggerConfigFromFile(cfg.LogConfigurationFile); err != nil {
			zap.S().Error(err)
		}
	}
	if newLogger == nil {
		if newLogger, err = config.GetStandartLogger(cfg.LogStandartChois); err != nil {
			zap.S().Fatal(err)
		}
	}
	if err = daemon.Run(cfg, newLogger); err != nil {
		newLogger.Error("Daemon run error", zap.Error(err))
	}

}
