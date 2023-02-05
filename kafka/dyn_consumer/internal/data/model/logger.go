package model

import "experiment_go/kafka/dyn_consumer/internal/pkg/logger"

func LoggerConfig(cfg AppConfig) logger.LogConfig {
	var logFileConfigs []logger.LogFileConfig
	for _, fileConf := range cfg.Logger.LogFileConfigs {
		logFileConfigs = append(logFileConfigs, logger.LogFileConfig{
			Levels:           fileConf.Levels,
			IsAccessLog:      fileConf.IsAccessLog,
			FullpathFilename: fileConf.FullpathFilename,
			MaxSize:          fileConf.MaxSize,
			MaxAge:           fileConf.MaxAge,
			MaxBackups:       fileConf.MaxBackups,
			LocalTime:        fileConf.LocalTime,
			Compress:         fileConf.Compress,
		})
	}

	return logger.LogConfig{
		EnableStdout:   cfg.Logger.EnableStdout,
		EnableLogFile:  cfg.Logger.EnableLogFile,
		CallerSkipSet:  cfg.Logger.CallerSkipSet,
		CallerSkip:     cfg.Logger.CallerSkip,
		LogFileConfigs: logFileConfigs,
	}
}
