package apphttp

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/esmailemami/eshop/consts"
	"github.com/esmailemami/eshop/errors"
	"github.com/esmailemami/eshop/services/logger"
)

func ErrorResponseHandler(ctx *HttpContext, err error) {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	switch e := err.(type) {
	case *errors.ValidationError:
		{
			logger.Default().WithField("StatusCode", strconv.Itoa(e.StatusCode)).Error(err.Error())
			data, err := json.Marshal(e)
			if err != nil {
				_ = ctx.QuickResponse(err.Error(), http.StatusInternalServerError)
				return
			}
			ctx.ResponseWriter.WriteHeader(e.StatusCode)
			_, _ = ctx.ResponseWriter.Write(data)

		}
	case *errors.RecordNotFoundError:
		{
			_ = ctx.QuickResponse(e.Error(), e.StatusCode)
			if e.LogError != nil {
				logger.Default().WithField("StatusCode", strconv.Itoa(e.StatusCode)).Error(e.LogError.Error())
			}
		}
	case *errors.BadRequestError:
		{
			_ = ctx.QuickResponse(e.Error(), e.StatusCode)
			if e.LogError != nil {
				logger.Default().WithField("StatusCode", strconv.Itoa(e.StatusCode)).Error(e.LogError.Error())
			}
		}
	case *errors.UnauthorizedError:
		{
			_ = ctx.QuickResponse(e.Error(), e.StatusCode)
			if e.LogError != nil {
				logger.Default().WithField("StatusCode", strconv.Itoa(e.StatusCode)).Error(e.LogError.Error())
			}
		}
	case *errors.ForbiddenError:
		{
			_ = ctx.QuickResponse(e.Error(), e.StatusCode)
			if e.LogError != nil {
				logger.Default().WithField("StatusCode", strconv.Itoa(e.StatusCode)).Error(e.LogError.Error())
			}
		}
	case *errors.InternalServerError:
		{
			_ = ctx.QuickResponse(e.Error(), e.StatusCode)
			if e.LogError != nil {
				logger.Default().WithField("StatusCode", strconv.Itoa(e.StatusCode)).Error(e.LogError.Error())
			}
		}
	default:
		{
			_ = ctx.QuickResponse(consts.InternalServerError, http.StatusInternalServerError)
			logger.Default().WithField("StatusCode", strconv.Itoa(http.StatusInternalServerError)).Error(e.Error())
		}
	}
}
