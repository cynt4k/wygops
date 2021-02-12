package router

import (
	"fmt"
	"net/http"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router/extension"
	"github.com/cynt4k/wygops/internal/router/middlewares"
	v1 "github.com/cynt4k/wygops/internal/router/v1"
	service "github.com/cynt4k/wygops/internal/services"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leandro-lugaresi/hub"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Router : Struct for the router class
type Router struct {
	e  *echo.Echo
	v1 *v1.Handlers
}

func Setup(
	hub *hub.Hub,
	db *gorm.DB,
	repo repository.Repository,
	config *config.Config,
	ss *service.Services,
	logger *zap.Logger,
) *echo.Echo {
	r := newRouter(hub, db, repo, config, ss, logger.Named("router"))

	r.e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, http.StatusText(http.StatusOK)) })
	r.e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	}, middlewares.AccessLoggingIgnore(), middlewares.MetricsIgnore())

	r.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()), middlewares.MetricsIgnore())

	r.v1.Setup(r.e.Group("/api"))

	if config.DevMode {
		routes := r.e.Routes()

		logger.Info("Registered routes are:")
		for _, route := range routes {
			logger.Info(fmt.Sprintf("Method: %s\tPath: %s\tFunction: %s", route.Method, route.Path, route.Name))
		}
	}

	return r.e
}

func newEcho(logger *zap.Logger, config *config.Config, repo repository.Repository) *echo.Echo {
	const maxAge = 300
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = extension.ErrorHandler(logger)

	e.Use(middlewares.RequestID())
	e.Use(middlewares.AccessLogging(logger.Named("access_log"), config.DevMode))
	e.Use(extension.Wrap(repo, *config))
	e.Use(middlewares.RequestCounter())
	e.Use(middlewares.HTTPDuration())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		ExposeHeaders: []string{echo.HeaderXRequestID},
		AllowHeaders:  []string{echo.HeaderContentType, echo.HeaderAuthorization},
		MaxAge:        maxAge,
	}))
	return e
}
