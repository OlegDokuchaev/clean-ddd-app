package admin

import (
	"api-gateway/internal/port/output/auth/admin"
)

type AuthImpl struct {
	config *Config
}

func NewAuth(config *Config) admin.Auth {
	return &AuthImpl{config: config}
}

func (a *AuthImpl) Validate(token string) bool {
	return a.config.Token == token
}

var _ admin.Auth = (*AuthImpl)(nil)
