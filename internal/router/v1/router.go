package v1

import (
	"net/http"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router/middlewares"
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/cynt4k/wygops/internal/services/wireguard"
	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

// Handlers : Handler struct
type Handlers struct {
	Repo      repository.Repository
	Bus       *hub.Hub
	Logger    *zap.Logger
	Config    *config.Config
	LDAP      ldap.LDAP
	Wireguard wireguard.Wireguard
}

// Config : Config struct
type Config struct {
	Version  string
	Revision string
}

func (h *Handlers) Setup(e *echo.Group) {
	apiNoAuth := e.Group("/v1")
	{
		apiNoAuth.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, http.StatusText(http.StatusOK)) })

		apiNoAuthAuth := apiNoAuth.Group("/auth")
		{
			apiNoAuthAuth.POST("/login", h.LoginUser)
			apiNoAuthAuth.POST("/logout", h.LogoutUser)
		}
	}

	api := e.Group("/v1", middlewares.CheckAuthentification(h.Config.API.JWT.Secret))
	{
		apiProfile := api.Group("/profile")
		{
			apiProfile.GET("", h.GetOwnProfile)
			apiProfile.PATCH("/protect", h.UserSetProtectPassword)
		}
	}
}

// Init : Initialize the v1 Routes
// func (h *Handlers) Init(g *gin.RouterGroup) {
// 	api := g.Group("/v1")
// 	api.GET("/", func(c *gin.Context) { c.String(http.StatusOK, http.StatusText(http.StatusOK)) })
// 	api.GET("/test", func(c *gin.Context) {
// 		c.String(http.StatusOK, http.StatusText(http.StatusOK))
// 	})
// 	api.GET("/test2", func(c *gin.Context) {
// 		ldapUser, err := h.LDAP.FindUser("developer", true)
// 		if err != nil {
// 			c.JSON(http.StatusNoContent, gin.H{
// 				"message": "no content",
// 			})
// 			return
// 		}

// 		user := models.User{
// 			Username:        ldapUser.Username,
// 			ProtectPassword: "asdf",
// 			FirstName:       ldapUser.FirstName,
// 			LastName:        ldapUser.LastName,
// 			Mail:            ldapUser.Mail,
// 			Type:            "ldap",
// 		}

// 		savedUser, err := h.Repo.CreateUser(&user)

// 		if err == repository.ErrAlreadyExists {
// 			savedUser, err = h.Repo.GetUserByUsername(user.Username)
// 			if err != nil {
// 				c.JSON(http.StatusNoContent, gin.H{
// 					"message": "no content",
// 				})
// 				return
// 			}
// 		}

// 		peer, err := h.Wireguard.CreatePeer()

// 		if err != nil {
// 			c.JSON(http.StatusNoContent, gin.H{
// 				"message": "no content",
// 			})
// 			return
// 		}

// 		device := models.Device{
// 			UserID:      savedUser.ID,
// 			IPv4Address: peer.IPV4Address.String(),
// 			IPv6Address: peer.IPV6Address.String(),
// 			PrivateKey:  peer.PrivateKey.String(),
// 			PublicKey:   peer.PublicKey.String(),
// 			Name:        "device",
// 		}

// 		deviceSaved, err := h.Repo.CreateDevice(&device)

// 		if err != nil {
// 			c.JSON(http.StatusNoContent, gin.H{
// 				"message": "no content",
// 			})
// 		}

// 		c.JSON(http.StatusOK, deviceSaved)
// 	})
// }
