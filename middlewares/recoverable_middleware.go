package middlewares

import (
	"errors"
	"runtime"
	"sync"

	"github.com/lab259/http"
)

var stackBuffPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4098)
	},
}

func RecoverableMiddleware(req http.Request, res http.Response, next http.Handler) http.Result {
	defer func() {
		if recoveryData := recover(); recoveryData != nil {
			logger := Logger(req)

			stack := stackBuffPool.Get().([]byte)
			defer stackBuffPool.Put(stack)

			n := runtime.Stack(stack, false)
			logger.Criticalf("panicked: %s", recoveryData)
			logger.Debug(string(stack[:n]))

			if err, ok := recoveryData.(error); ok {
				res.Error(err)
			} else {
				res.Error(errors.New("unexpected panic"))
			}
		}
	}()
	return next(req, res)
}
