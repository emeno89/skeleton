package security

import (
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

const (
	jwtAudience = "users"
)

type claims struct {
	jwt.RegisteredClaims
}

type JwtGen struct {
	secret []byte
}

func newJwtGen(cfg jwtConfig) *JwtGen {
	return &JwtGen{
		secret: []byte(cfg.Secret),
	}
}

func (s *JwtGen) ParseUserId(token string) (userId string, err error) {
	c := claims{}

	sFn := func(token *jwt.Token) (interface{}, error) { return s.secret, nil }

	if _, err := jwt.ParseWithClaims(strings.TrimPrefix(token, "Bearer "), &c, sFn, jwt.WithAudience(jwtAudience)); err != nil {
		return "", err
	}

	return c.Subject, nil
}
