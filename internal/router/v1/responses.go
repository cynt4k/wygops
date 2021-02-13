package v1

import (
	"github.com/cynt4k/wygops/internal/models"
	"github.com/cynt4k/wygops/internal/router/consts"
	"github.com/cynt4k/wygops/internal/router/extension"
)

type UserLogin struct {
	Username  string `json:"username"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type OwnProfile struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Mail        string `json:"mail"`
	Initialized bool   `json:"initialized"`
}

func formatOwnProfile(user *models.User) *extension.Response {
	initialzied := false

	if user.ProtectPassword != "" {
		initialzied = true
	}
	return &extension.Response{
		Message: string(consts.I18nResponseOK),
		Data: OwnProfile{
			Username:    user.Username,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Mail:        user.Mail,
			Initialized: initialzied,
		},
	}
}
