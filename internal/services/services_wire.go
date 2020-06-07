// +build wireinject

package service

import (
	"github.com/google/wire"
)

// ProviderSet : Services dependency injection
var ProviderSet = wire.NewSet(wire.FieldsOf(new(*Services),
	"User",
	"Ldap",
))
