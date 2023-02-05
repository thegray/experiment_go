package logger

import "time"

type LogConfig struct {
	EnableStdout   bool
	EnableLogFile  bool
	EnableELK      bool
	CallerSkipSet  bool
	CallerSkip     int
	LogFileConfigs []LogFileConfig
	ELKConfig      ELKConfig
}

type LogFileConfig struct {
	Levels []string
	// When writing this code, zap doesn't have a different log level for Access and
	// zapcore.Core will write logs to a WriteSyncer based on the zap level.
	// If we want to add a new level of logs, we have to extend the zap library.
	// So as a workaround for this problem, IsAccessLog field is introduced to "force" the logger
	// to choose whether it will write info level log to access.log.
	IsAccessLog      bool
	FullpathFilename string
	MaxSize          int
	MaxAge           int
	MaxBackups       int
	LocalTime        bool
	Compress         bool
}

type ELKConfig struct {
	Host  string
	Index string

	Username       string
	Password       string
	TLSCertificate string

	// Size specifies the maximum amount of data the writer will buffered
	// before flushing.
	//
	// Defaults to 256 kB if unspecified.
	BufferSize int

	// FlushInterval specifies how often the writer should flush data if
	// there have been no writes.
	//
	// Defaults to 30 seconds if unspecified.
	FlushInterval time.Duration
}
