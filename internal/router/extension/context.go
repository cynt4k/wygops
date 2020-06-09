package extension

import (
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router/consts"
	"github.com/gin-gonic/gin"
)

type Context struct {
	gin.Context
}

// Wrap : Wrap the repository into the request
func Wrap(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.KeyRepo, repo)
		c.Next()
	}
}
