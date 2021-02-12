package v1

import (
	"net/http"

	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/labstack/echo/v4"
)

type PostLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) LoginUser(c echo.Context) error {
	var req PostLoginRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	authSuccess, err := h.LDAP.CheckLogin(req.Username, req.Password)

	if err != nil || !authSuccess {
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "user login failed")
	}

	user, err := h.Repo.GetUserByUsername(req.Username)

	if err != nil {
		switch err {
		case repository.ErrNotFound:
			userLdap, err := h.LDAP.FindUser(req.Username, true)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			user = &models.User{
				Username:  userLdap.Username,
				Mail:      userLdap.Mail,
				FirstName: userLdap.FirstName,
				LastName:  userLdap.LastName,
				Type:      "ldap",
			}

			user, err = h.Repo.CreateUser(user)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	token, err := createJWTToken(user, h.Config.API.JWT.Secret, h.Config.API.JWT.LifeTime)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, UserLogin{
		Username:  user.Username,
		Token:     token.AccessToken,
		ExpiresAt: token.Expires,
	})
}

// TODO: Implement logout
func (h *Handlers) LogoutUser(c echo.Context) error {
	return nil
}
