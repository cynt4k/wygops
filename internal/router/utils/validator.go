package utils

import (
	"context"

	"github.com/cynt4k/wygops/internal/router/consts"
	"github.com/labstack/echo/v4"
)

type ctxKey int

const (
	repoCtxKey ctxKey = iota
)

func NewRequestValidateContext(c echo.Context) context.Context {
	return context.WithValue(context.Background(), repoCtxKey, c.Get(consts.KeyRepo))
}
