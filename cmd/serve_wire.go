// +build wireinject

package cmd

import (
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router"
	service "github.com/cynt4k/wygops/internal/services"
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/cynt4k/wygops/internal/services/user"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

func newHttpServer(hub *hub.Hub, db *gorm.DB, repo repository.Repository, logger *zap.Logger) (*HTTPServer, error) {
	wire.Build(
		router.Init,
		ldap.NewService,
		user.NewService,
		wire.Struct(new(service.Services), "*"),
		wire.Struct(new(HTTPServer), "*"),
	)
	return nil, nil
}
