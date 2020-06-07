// +build wireinject

package router

import (
	"github.com/cynt4k/wygops/internal/repository"
	v1 "github.com/cynt4k/wygops/internal/router/v1"
	service "github.com/cynt4k/wygops/internal/services"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

func newRouter(hub *hub.Hub, db *gorm.DB, repo repository.Repository, ss *service.Services, logger *zap.Logger) *Router {
	wire.Build(
		service.ProviderSet,
		newGin,
		wire.Struct(new(v1.Handlers), "*"),
	)
	return nil
}
