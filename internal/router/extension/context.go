package extension

import (
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router/consts"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type response struct {
	Code int `json:"code"`
	Response
}

type Context struct {
	echo.Context
}

func (c *Context) JSON(code int, i interface{}) (err error) {
	switch res := i.(type) {
	case Response:
		return c.Context.JSON(code, &response{
			Code:     code,
			Response: res,
		})
	default:
		return c.Context.JSON(code, &response{
			Code: code,
			Response: Response{
				Message: string(consts.I18nResponseOK),
				Data:    i,
			},
		})
	}
}

func json(c echo.Context, code int, i interface{}, cfg jsoniter.StreamPool) error {
	// use as an jsoniter.API
	stream := cfg.BorrowStream(c.Response())
	defer cfg.ReturnStream(stream)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(code)
	stream.WriteVal(i)
	stream.WriteRaw("\n")
	return stream.Flush()
}

func Wrap(repo repository.Repository) echo.MiddlewareFunc {
	return func(n echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(consts.KeyRepo, repo)
			return n(&Context{Context: c})
		}
	}
}
