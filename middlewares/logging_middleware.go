package middlewares

import (
	"github.com/lab259/http"
	"github.com/lab259/rlog"
)

func Logger(req http.Request) rlog.Logger {
	return rlog.WithField("request_id", req.Raw().ID())
}

func LoggingMiddleware(req http.Request, res http.Response, next http.Handler) http.Result {
	logger := Logger(req)
	logger.Infof("%s: %s", req.Method(), req.Path())
	return next(req, res)
}
