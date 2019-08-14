package middlewares

import (
	"github.com/lab259/hermes"
	"github.com/lab259/rlog/v2"
)

func Logger(req hermes.Request) rlog.Logger {
	return rlog.WithField("request_id", req.Raw().ID())
}

func LoggingMiddleware(req hermes.Request, res hermes.Response, next hermes.Handler) hermes.Result {
	logger := Logger(req)
	logger.Infof("%8s: %s", req.Method(), req.Path())
	return next(req, res)
}
