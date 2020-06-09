package router

import (
	"net/http"
	"time"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router/extension"
	v1 "github.com/cynt4k/wygops/internal/router/v1"
	service "github.com/cynt4k/wygops/internal/services"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	app = Router{}
)

// Router : Struct for the router class
type Router struct {
	app *gin.Engine
	v1  *v1.Handlers
}

// Init : Initialize the Router
func Init(hub *hub.Hub, db *gorm.DB, repo repository.Repository, ss *service.Services, logger *zap.Logger) *gin.Engine {
	r := newRouter(hub, db, repo, ss, logger.Named("router"))
	api := r.app.Group("/api")

	r.app.GET("/", func(c *gin.Context) { c.String(http.StatusOK, http.StatusText(http.StatusOK)) })
	api.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, http.StatusText(http.StatusOK)) })
	api.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.v1.Init(api)

	return r.app
}

func newGin(logger *zap.Logger, repo repository.Repository) *gin.Engine {
	if !config.GetConfig().DevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	g.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	g.Use(ginzap.RecoveryWithZap(logger, true))
	g.Use(extension.Wrap(repo))

	return g
}
