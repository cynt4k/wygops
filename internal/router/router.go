package router

import (
	"time"

	"github.com/cynt4k/wygops/internal/repository"
	v1 "github.com/cynt4k/wygops/internal/router/v1"
	service "github.com/cynt4k/wygops/internal/services"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
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
func Init(hub *hub.Hub, db *gorm.DB, repo repository.Repository, ss *service.Services, logger *zap.Logger) error {
	r := newRouter(hub, db, repo, ss, logger.Named("router"))
	r.app = gin.New()

	return nil
}

func newGin(logger *zap.Logger, repo *repository.Repository) *gin.Engine {
	g := gin.New()
	g.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	g.Use(ginzap.RecoveryWithZap(logger, true))

	return g
}
