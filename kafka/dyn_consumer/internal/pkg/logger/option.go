package logger

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type config struct {
	ws            []zapcore.WriteSyncer
	callerSkipSet bool
	callerSkip    int
}

type Option func(conf *config)

func AddWriter(w io.Writer, concurrentSafe bool) Option {
	return func(conf *config) {
		ws := zapcore.AddSync(w)
		if !concurrentSafe {
			ws = zapcore.Lock(ws)
		}
		conf.ws = append(conf.ws, ws)
	}
}

func WithCaller(callerSkip int) Option {
	return func(conf *config) {
		// Add 1 skip because the callstack will increase 1 from spid/logger log method
		conf.callerSkip = callerSkip + 1
		conf.callerSkipSet = true
	}
}
