package v1

import (
	"net/http"

	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/gin-gonic/gin"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

// Handlers : Handler struct
type Handlers struct {
	Repo   repository.Repository
	Bus    *hub.Hub
	Logger *zap.Logger
	User   *ldap.LDAP
}

// Config : Config struct
type Config struct {
	Version  string
	Revision string
}

// Init : Initialize the v1 Routes
func (h *Handlers) Init(g *gin.RouterGroup) {
	api := g.Group("/v1")
	api.GET("/", func(c *gin.Context) { c.String(http.StatusOK, http.StatusText(http.StatusOK)) })
}
