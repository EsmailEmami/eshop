package apphttp

import (
	"net/http"
	"runtime/debug"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/services/logger"
)

func Handler(
	f func(*HttpContext) error,
	middlewares ...func(ctx *HttpContext) error,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fields := make(map[string]string)
		fields["remote"] = r.RemoteAddr
		fields["method"] = r.Method
		fields["url"] = r.URL.String()

		logger.Default().WithFields(fields).Info("request")

		defer func() {
			if r := recover(); r != nil {
				logger.Default().Errorf("%v: %s", r, debug.Stack())
				http.Error(w, consts.InternalServerError, http.StatusInternalServerError)
			}
		}()

		ctx := NewHttpContext(w, r)

		for _, m := range middlewares {
			middlewareErr := m(ctx)
			if middlewareErr != nil {
				ErrorResponseHandler(ctx, middlewareErr)
				return
			}
		}

		err := f(ctx)
		if err != nil {
			ErrorResponseHandler(ctx, err)
		}
	})
}
