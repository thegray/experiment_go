package logger

import (
	"context"
	"errors"
	"experiment_go/kafka/dyn_consumer/internal/pkg/logger/contextid"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"moul.io/zapfilter"
)

// DefaultLogger Skip 1 Caller because default logger method will call log method from the Logger struct
var DefaultLogger = New(AddWriter(os.Stdout, false), WithCaller(1))

// Logger wrap underlying logger library
type Logger struct {
	logger    *zap.SugaredLogger
	threshold LogLevel
	stopFn    func()
}

func (l *Logger) Stop() {
	if l.stopFn == nil {
		return
	}

	l.stopFn()
}

type LogLevel int8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	OFF
)

const (
	maskedStr = "[Masked]"
)

func (l *Logger) SetThreshold(level LogLevel) {
	l.threshold = level
}

// Set the global DefaultLogger to l
func SetDefaultLogger(l Logger) {
	DefaultLogger = l
}

// New instantiates new logger
func New(opts ...Option) Logger {
	var conf config
	for _, opt := range opts {
		opt(&conf)
	}

	jsonEncoder := zapJSONEncoder()
	level := zapLevel()
	var cores []zapcore.Core
	for _, ws := range conf.ws {
		cores = append(cores, zapcore.NewCore(jsonEncoder, ws, level))
	}

	core := zapcore.NewTee(cores...)
	core = zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
	var zapOpts []zap.Option
	if conf.callerSkipSet {
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddCallerSkip(conf.callerSkip))
	}
	logger := zap.New(core, zapOpts...)
	defer logger.Sync()
	L := logger.Sugar()

	return Logger{logger: L, threshold: INFO}
}

// Instantiates new logger based on config supplied by user
//
// When enabling `EnableELK`, another goroutine is spawned to flush the buffer periodically.
// Logger's Stop() function must be deferred call after calling this function.
func NewFromConfig(conf LogConfig) (*Logger, error) {
	var cores []zapcore.Core

	if !conf.EnableLogFile && !conf.EnableStdout && !conf.EnableELK {
		return nil, errors.New("invalid configuration, must enable stdout or logfile or elk")
	}

	jsonEncoder := zapJSONEncoder()

	if conf.EnableLogFile {
		if len(conf.LogFileConfigs) == 0 {
			return nil, errors.New("log file configurations must not be empty")
		}

		for _, logFileConfig := range conf.LogFileConfigs {
			core, err := createFileHandlerCore(jsonEncoder, logFileConfig)
			if err != nil {
				return nil, err
			}

			cores = append(cores, core)
		}
	}

	if conf.EnableStdout {
		core := createStdoutHanlderCore(jsonEncoder)
		cores = append(cores, core)
	}

	var stopFn func()

	L := createZapLogger(cores, conf.CallerSkipSet, conf.CallerSkip)

	return &Logger{logger: L, threshold: INFO, stopFn: stopFn}, nil
}

// Create zap core that specifically handle writing log to file
func createFileHandlerCore(encoder zapcore.Encoder, logFileConfig LogFileConfig) (zapcore.Core, error) {
	writer := &lumberjack.Logger{
		Filename:   logFileConfig.FullpathFilename,
		MaxSize:    logFileConfig.MaxSize,
		MaxBackups: logFileConfig.MaxBackups,
		MaxAge:     logFileConfig.MaxAge,
		LocalTime:  logFileConfig.LocalTime,
		Compress:   logFileConfig.Compress,
	}
	writeSyncer := zapcore.Lock(zapcore.AddSync(writer))

	core := zapcore.NewCore(encoder, writeSyncer, zapLevel())
	levelFilterPattern := strings.Join(logFileConfig.Levels, ",")

	var filter string = fmt.Sprintf("%s:%s", levelFilterPattern, "*,-access*")
	if logFileConfig.IsAccessLog {
		filter = fmt.Sprintf("%s:%s", levelFilterPattern, "access")
	}

	levelFilterFunction, err := zapfilter.ParseRules(filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", "invalid log file levels configuration", err.Error())
	}
	filteredCore := zapfilter.NewFilteringCore(core, levelFilterFunction)

	return filteredCore, nil
}

// Create zap core that specifically handle writing log to stdout
func createStdoutHanlderCore(encoder zapcore.Encoder) zapcore.Core {
	writer := os.Stdout
	writeSyncer := zapcore.AddSync(writer)
	return zapcore.NewCore(encoder, writeSyncer, zapLevel())
}

// Create zap logger based on slice of zapcore.Core
func createZapLogger(cores []zapcore.Core, callerSkipSet bool, callerSkip int) *zap.SugaredLogger {
	core := zapcore.NewTee(cores...)
	core = zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)

	var zapOpts []zap.Option
	if callerSkipSet {
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddCallerSkip(callerSkip))
	}

	logger := zap.New(core, zapOpts...)
	defer logger.Sync()

	return logger.Sugar()
}

// Debug using the default logger to log the message on debug level with additional key value when provided
func Debug(msg string, kv ...interface{}) {
	DefaultLogger.Debug(msg, kv...)
}

// Info using the default logger to log the message on info level with additional key value when provided
func Info(msg string, kv ...interface{}) {
	DefaultLogger.Info(msg, kv...)
}

// Access using the default logger to log the message on info level with additional key value when provided
// This function is created to differentiate the log used for access logging purpose
// When the logs are written to file it should be written to different log file (e.g. Info to data.log, while Access to access.log)
func Access(msg string, kv ...interface{}) {
	DefaultLogger.Access(msg, kv...)
}

// Warn using the default logger to log the message on warn level with additional key value when provided
func Warn(msg string, kv ...interface{}) {
	DefaultLogger.Warn(msg, kv...)
}

// Error using the default logger to log the message on error level with the error detail and additional key value when provided
// inline causes stacktrace to mistakenly print out this function instead of the caller.
//go:noinline
func Error(msg string, kv ...interface{}) {
	DefaultLogger.Error(msg, kv...)
}

// DebugCtx using the default logger to log the message on debug level with additional key value when provided
func DebugCtx(ctx context.Context, msg string, kv ...interface{}) {
	DefaultLogger.DebugCtx(ctx, msg, kv...)
}

// InfoCtx using the default logger to log the message on info level with additional key value when provided
func InfoCtx(ctx context.Context, msg string, kv ...interface{}) {
	DefaultLogger.InfoCtx(ctx, msg, kv...)
}

// Same with Access, with additional key value for context
func AccessCtx(ctx context.Context, msg string, kv ...interface{}) {
	DefaultLogger.AccessCtx(ctx, msg, kv...)
}

// WarnCtx using the default logger to log the message on warn level with additional key value when provided
func WarnCtx(ctx context.Context, msg string, kv ...interface{}) {
	DefaultLogger.WarnCtx(ctx, msg, kv...)
}

// ErrorCtx using the default logger to log the message on error level with the error detail and additional key value when provided
// inline causes stacktrace to mistakenly print out this function instead of the caller.
//go:noinline
func ErrorCtx(ctx context.Context, msg string, kv ...interface{}) {
	DefaultLogger.ErrorCtx(ctx, msg, kv...)
}

// Deprecated: Timestamp will be set by default
// WithTime add log timestamp to every log produced
func (l Logger) WithTime() Logger {
	return l
}

// With add additional information to every log produced
func (l Logger) With(key string, value interface{}) Logger {
	value = maskInterface(value, false)
	l.logger = l.logger.With(key, value)
	return l
}

// Debug log the message on debug level with additional key value when provided
func (l Logger) Debug(msg string, kv ...interface{}) {
	if DEBUG < l.threshold {
		return
	}
	l.logger.Debugw(msg, mask(kv...)...)
}

// Info log the message on info level with additional key value when provided
func (l Logger) Info(msg string, kv ...interface{}) {
	if INFO < l.threshold {
		return
	}
	l.logger.Infow(msg, mask(kv...)...)
}

// Access log the message on info level with additional key value when provided
// See also: additional comment for DefaultLogger Access function
func (l Logger) Access(msg string, kv ...interface{}) {
	if INFO < l.threshold {
		return
	}
	// a namespace called "access" is added to this method
	// so that the logger that calls this function can filter
	// log with INFO severity should be treated as an access log or data log
	l.logger.Named("access").Infow(msg, mask(kv...)...)
}

// Warn log the message on warn level with additional key value when provided
func (l Logger) Warn(msg string, kv ...interface{}) {
	if WARN < l.threshold {
		return
	}
	l.logger.Warnw(msg, mask(kv...)...)
}

// Error log the message on error level with the error detail and additional key value when provided
func (l Logger) Error(msg string, kv ...interface{}) {
	if ERROR < l.threshold {
		return
	}

	l.logger.Errorw(msg, mask(kv...)...)
}

// DebugCtx log the message on debug level with additional key value when provided
func (l Logger) DebugCtx(ctx context.Context, msg string, kv ...interface{}) {
	if DEBUG < l.threshold {
		return
	}

	contextID := contextid.Value(ctx)
	if contextID != "" {
		kv = append(kv, "context_id", contextID)
	}

	l.logger.Debugw(msg, mask(kv...)...)
}

// InfoCtx log the message on info level with additional key value when provided
func (l Logger) InfoCtx(ctx context.Context, msg string, kv ...interface{}) {
	if INFO < l.threshold {
		return
	}

	contextID := contextid.Value(ctx)
	if contextID != "" {
		kv = append(kv, "context_id", contextID)
	}

	l.logger.Infow(msg, mask(kv...)...)
}

// AccessCtx log the message on info level with additional key value when provided
func (l Logger) AccessCtx(ctx context.Context, msg string, kv ...interface{}) {
	if INFO < l.threshold {
		return
	}

	contextID := contextid.Value(ctx)
	if contextID != "" {
		kv = append(kv, "context_id", contextID)
	}

	l.logger.Named("access").Infow(msg, mask(kv...)...)
}

// WarnCtx log the message on warn level with additional key value when provided
func (l Logger) WarnCtx(ctx context.Context, msg string, kv ...interface{}) {
	if WARN < l.threshold {
		return
	}

	contextID := contextid.Value(ctx)
	if contextID != "" {
		kv = append(kv, "context_id", contextID)
	}

	l.logger.Warnw(msg, mask(kv...)...)
}

// ErrorCtx log the message on error level with the error detail and additional key value when provided
func (l Logger) ErrorCtx(ctx context.Context, msg string, kv ...interface{}) {
	if ERROR < l.threshold {
		return
	}

	contextID := contextid.Value(ctx)
	if contextID != "" {
		kv = append(kv, "context_id", contextID)
	}

	l.logger.Errorw(msg, mask(kv...)...)
}

func mask(kv ...interface{}) []interface{} {
	n := len(kv)
	var params []interface{}
	for i := 0; i < n-1; i += 2 {
		if _, ok := kv[i].(string); ok {
			params = append(params, kv[i], maskInterface(kv[i+1], false))
		}
	}
	return params
}

func convert(v interface{}) interface{} {
	switch bVal := v.(type) {
	case fmt.Stringer:
		v = bVal.String()
	case []byte:
		// when we receive []byte, convert it to string to prevent logger from printing the ascii values.
		// non-printable characters will be escaped (e.g \u007f)
		v = string(bVal)
	}
	return v
}
