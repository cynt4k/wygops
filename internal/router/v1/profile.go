package v1

import (
	"errors"
	"net/http"

	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/router/extension/herror"
	"github.com/cynt4k/wygops/pkg/util/optional"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetOwnProfile(c echo.Context) error {
	userToken, err := extractUserInfos(c)

	if err != nil {
		return err
	}

	user, err := h.Repo.GetUserByUsername(userToken.Username)

	if err != nil {
		return herror.InternalServerError(err)
	}

	return c.JSON(http.StatusOK, formatOwnProfile(user))
}

type PostUserSetProtectPassword struct {
	ProtectPassword string `json:"protectPassword"`
}

func (h *Handlers) UserSetProtectPassword(c echo.Context) error {
	var req PostUserSetProtectPassword
	userToken, err := extractUserInfos(c)

	if err != nil {
		return err
	}

	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.Repo.GetUserByUsername(userToken.Username)

	if err != nil {
		return herror.InternalServerError(err)
	}

	if user.ProtectPassword != "" && user.Cipher != "" {
		return herror.BadRequest("profile already initialized")
	}
	if user.ProtectPassword != "" && user.Cipher == "" {
		return herror.InternalServerError(errors.New("profile has ciper but empty protect password - ask admin"))
	}
	if user.ProtectPassword == "" && user.Cipher != "" {
		return herror.InternalServerError(errors.New("profile has protect password but no cipher - ask admin"))
	}
	user, err = h.Repo.UpdateUser(user.ID, repository.UpdateUserArgs{
		ProtectPassword: optional.NewString(req.ProtectPassword, true),
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, formatOwnProfile(user))
}
