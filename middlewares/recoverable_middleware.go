package middlewares

import (
	"errors"
	"runtime"
	"sync"

	"github.com/lab259/hermes"
)

var stackBuffPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4098)
	},
}

func RecoverableMiddleware(req hermes.Request, res hermes.Response, next hermes.Handler) (r hermes.Result) {
	defer func() {
		if recoveryData := recover(); recoveryData != nil {
			logger := Logger(req)

			stack := stackBuffPool.Get().([]byte)
			defer stackBuffPool.Put(stack)

			n := runtime.Stack(stack, false)
			logger.Criticalf("panicked: %s", recoveryData)
			logger.Debug(string(stack[:n]))

			if err, ok := recoveryData.(error); ok {
				r = res.Error(err)
			} else {
				r = res.Error(errors.New("unexpected panic"))
			}
		}
	}()
	return next(req, res)
}
