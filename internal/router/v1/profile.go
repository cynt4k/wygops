package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) GetOwnProfile(c echo.Context) error {
	user, err := extractUserInfos(c)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}
