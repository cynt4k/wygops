package v1

import (
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

type Handlers struct {
	Repo   repository.Repository
	Bus    *hub.Hub
	Logger *zap.Logger
}

type Config struct {
	Version  string
	Revision string
}

// func (h *Handlers) Init(group *gin.RouterGroup) error {

// }

func (h *Handlers) Init(g *gin.RouterGroup) {
	api := g.Group("/v1")
	api.Use()
}
