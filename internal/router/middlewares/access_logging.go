package middlewares

import (
	"strconv"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/cynt4k/wygops/internal/router/consts"
	"github.com/cynt4k/wygops/internal/router/extension"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func AccessLoggingIgnore() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(consts.KeyIgnoreLogging, true)
			return next(c)
		}
	}
}

func AccessLogging(logger *zap.Logger, dev bool) echo.MiddlewareFunc {
	if dev {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				start := time.Now()
				if err := next(c); err != nil {
					c.Error(err)
				}
				stop := time.Now()

				ignoreLogging := c.Get(consts.KeyIgnoreLogging)
				if ignoreLogging != nil {
					if ignoreLogging.(bool) {
						return nil
					}
				}

				req := c.Request()
				res := c.Response()
				logger.Sugar().Infof("%3d | %s | %s %s %d", res.Status, stop.Sub(start), req.Method, req.URL, res.Size)
				return nil
			}
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			ignoreLogging := c.Get(consts.KeyIgnoreLogging)
			if ignoreLogging != nil {
				if ignoreLogging.(bool) {
					return nil
				}
			}

			req := c.Request()
			res := c.Response()
			logger.Info("",
				zap.String("requestId", extension.GetRequestID(c)),
				zapdriver.HTTP(&zapdriver.HTTPPayload{
					RequestMethod: req.Method,
					Status:        res.Status,
					UserAgent:     req.UserAgent(),
					RemoteIP:      c.RealIP(),
					Referer:       req.Referer(),
					Protocol:      req.Proto,
					RequestURL:    req.URL.String(),
					RequestSize:   req.Header.Get(echo.HeaderContentLength),
					ResponseSize:  strconv.FormatInt(res.Size, 10),
					Latency:       strconv.FormatFloat(stop.Sub(start).Seconds(), 'f', 9, 64) + "s",
				}))
			return nil
		}
	}
}
