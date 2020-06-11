package service

import (
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/cynt4k/wygops/internal/services/user"
	"github.com/cynt4k/wygops/internal/services/wireguard"
)

// Services : Services which exists
type Services struct {
	User      *user.Service
	Ldap      ldap.LDAP
	Wireguard wireguard.Wireguard
}
