// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router"
	"github.com/cynt4k/wygops/internal/services"
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/cynt4k/wygops/internal/services/user"
	"github.com/jinzhu/gorm"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from serve_wire.go:

func newHttpServer(hub2 *hub.Hub, db *gorm.DB, repo repository.Repository, logger *zap.Logger, config2 *config.ProviderLdap) (*HTTPServer, error) {
	userService := user.NewService(repo, hub2, logger)
	ldapLDAP, err := ldap.NewService(repo, config2)
	if err != nil {
		return nil, err
	}
	services := &service.Services{
		User: userService,
		Ldap: ldapLDAP,
	}
	engine := router.Init(hub2, db, repo, services, logger)
	httpServer := &HTTPServer{
		Logger: logger,
		SS:     services,
		Router: engine,
		Hub:    hub2,
		Repo:   repo,
	}
	return httpServer, nil
}
