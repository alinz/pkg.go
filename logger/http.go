package logger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

type HttpOptions struct {
	ExcludePing bool
	PingPath    string
}

func RequestLogger(opt HttpOptions) func(next http.Handler) http.Handler {
	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		var logEvent *zerolog.Event

		if opt.ExcludePing && status == http.StatusOK && r.URL.String() == opt.PingPath {
			return
		}

		if status >= 400 {
			logEvent = log.Logger.Error()
		} else {
			logEvent = log.Logger.Info()
		}

		logEvent.
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})
}
