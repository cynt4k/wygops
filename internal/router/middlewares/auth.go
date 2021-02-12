package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cynt4k/wygops/internal/router/consts"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func CheckAuthentification(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var baererToken string
			request := c.Request()
			baererTokenRaw := request.Header.Get(consts.HeaderAuthorization)

			if baererTokenRaw == "" {
				return echo.NewHTTPError(http.StatusForbidden, errors.New("no baerer token provided - refused"))
			}
			const tokenSplitSize = 2
			if tokenSplit := strings.Split(baererTokenRaw, " "); len(tokenSplit) == tokenSplitSize {
				baererToken = tokenSplit[1]
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid token format"))
			}

			tkn, err := jwt.Parse(baererToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexcpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, err)
			}

			if _, ok := tkn.Claims.(jwt.Claims); !ok && !tkn.Valid {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("token is invalid"))
			}

			c.Set(consts.KeyToken, tkn)
			return next(c)
		}
	}
}
