package v1

import (
	"fmt"
	"time"

	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/internal/router/consts"
	"github.com/cynt4k/wygops/internal/router/extension"
	"github.com/dgrijalva/jwt-go"
	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type JWTUser struct {
	Username     string `json:"username"`
	AuthStrategy string `json:"authStrategy"`
}

type JWTTokenDetails struct {
	AccessToken string
	Expires     int64
}

func extractUserInfos(c echo.Context) (*JWTUser, error) {
	switch tokenFormat := c.Get(consts.KeyToken).(type) {
	case *jwt.Token:
	default:
		return nil, fmt.Errorf("provided context token invalid format %T", tokenFormat)
	}
	token := c.Get(consts.KeyToken).(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		username, ok := claims["username"].(string)

		if !ok {
			return nil, fmt.Errorf("invalid token format - username attribute missing")
		}

		authStrategy, ok := claims["authStrategy"].(string)

		if !ok {
			return nil, fmt.Errorf("invalid token format - authStrategy attribute missing")
		}

		return &JWTUser{
			Username:     username,
			AuthStrategy: authStrategy,
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func createJWTToken(user *models.User, jwtSecret string, lifetime string) (*JWTTokenDetails, error) {
	duration, err := time.ParseDuration(lifetime)

	if err != nil {
		return nil, err
	}

	td := &JWTTokenDetails{
		Expires: time.Now().Add(duration).Unix(),
	}

	tokenClaims := jwt.MapClaims{}

	tokenClaims["username"] = user.Username
	tokenClaims["exp"] = td.Expires
	tokenClaims["authStrategy"] = user.Type

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	td.AccessToken, err = token.SignedString([]byte(jwtSecret))

	if err != nil {
		return nil, err
	}
	return td, nil
}

func bindAndValidate(c echo.Context, i interface{}, rules ...vd.Rule) error {
	return extension.BindAndValidate(c, i, rules...)
}
