package apphttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/esmailemami/eshop/apphttp/httpmodels"
	"github.com/esmailemami/eshop/consts"
	"github.com/go-chi/chi/v5"
)

type TransformableObject interface {
	ApplyDataTransformation()
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		Request:        r,
		ResponseWriter: w,
	}
}

type HttpContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	pathParams     map[string]string
}

func (ctx *HttpContext) QuickResponse(message string, statusCode int) error {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")

	res := httpmodels.SuccessResponse{Message: message}
	return ctx.JSON(res, statusCode)
}

func (ctx *HttpContext) JSON(object interface{}, statusCode int) error {
	ctx.ResponseWriter.WriteHeader(statusCode)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(ctx.ResponseWriter).Encode(&object)
}

func (ctx *HttpContext) Bind(object TransformableObject) error {
	if ctx.Request.Body == nil {
		return errors.New(consts.BadRequest)
	}
	err := json.NewDecoder(ctx.Request.Body).Decode(&object)
	if err != nil {
		return err
	}

	object.ApplyDataTransformation()
	return nil
}

// bind request data to given object without applying any kind of transformation to sanitizing data
func (ctx *HttpContext) BlindBind(object any) error {
	if ctx.Request.Body == nil {
		return errors.New(consts.BadRequest)
	}
	err := json.NewDecoder(ctx.Request.Body).Decode(&object)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *HttpContext) GetPathParam(name string) string {
	val := chi.URLParam(ctx.Request, name)
	return val
}

func (ctx *HttpContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.Request.Context().Deadline()
}

func (ctx *HttpContext) Done() <-chan struct{} {
	return ctx.Request.Context().Done()
}

func (ctx *HttpContext) Err() error {
	return ctx.Request.Context().Err()
}

func (ctx *HttpContext) Value(key interface{}) interface{} {
	return ctx.Request.Context().Value(key)
}

func (ctx *HttpContext) ClientIP() string {
	clientIP := ctx.Request.Header.Get("X-Forwarded-For")

	if clientIP == "" {
		clientIP = ctx.Request.RemoteAddr
	} else {
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	}

	return clientIP
}

func (ctx *HttpContext) UserAgent() string {
	return ctx.Request.UserAgent()
}

func (c *HttpContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.ResponseWriter, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 0,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *HttpContext) Get(key string) (any, bool) {
	val := c.Value(key)

	if val == nil {
		return nil, false
	}

	return val, true
}
