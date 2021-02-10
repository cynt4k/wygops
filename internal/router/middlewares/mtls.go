package middlewares

import (
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// func getAuthorication(peerCertificates []*x509.Certificate) (bool, error) {
// 	var authorizedKeys []interface{}
// }

func MTLSAuth(caPool *x509.CertPool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			certs := c.Request().TLS.PeerCertificates

			if len(certs) == 0 {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf("you are not permitted to request to '%s' - no client cert provided", c.Request().URL.Path),
				)
			}

			opts := x509.VerifyOptions{
				Roots:         caPool,
				Intermediates: x509.NewCertPool(),
				KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}

			for _, cert := range certs[1:] {
				opts.Intermediates.AddCert(cert)
			}

			_, err := certs[0].Verify(opts)

			if err != nil {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					fmt.Sprintf("you are not permitted to request to '%s' - %s", c.Request().URL.Path, err),
				)
			}

			return next(c)
		}
	}
}
