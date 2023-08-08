package jwtx

import (
	"context"
	"time"

	"github.com/ytf606/golibs/errorx"
	"github.com/ytf606/golibs/logx"
	"github.com/golang-jwt/jwt/v4"
)

type (
	Claims         = jwt.Claims
	StandardClaims = jwt.StandardClaims
	MapClaims      = jwt.MapClaims
)

type Jwter interface {
	Create(ctx context.Context, claims jwt.Claims) (string, error)
	Parse(ctx context.Context, tokenStr string, claims jwt.Claims) (*jwt.Token, error)
	SetExpiresAt(value int64) int64
	SetSignKey(value string) string
	SetIssuer(value string) string
}

type jwter struct {
	SignKey   string
	ExpiresAt int64
	Issuer    string
}

func NewJwt(signKey string) Jwter {
	return &jwter{
		SignKey: signKey,
	}
}

func (j *jwter) GetSignKey() string {
	return j.SignKey
}

func (j *jwter) SetSignKey(value string) string {
	j.SignKey = value
	return j.SignKey
}

func (j *jwter) SetExpiresAt(value int64) int64 {
	j.ExpiresAt = value
	return j.ExpiresAt
}

func (j *jwter) SetIssuer(value string) string {
	j.Issuer = value
	return j.Issuer
}

func (j *jwter) Create(ctx context.Context, claims jwt.Claims) (string, error) {
	tag := "[jwtx_jwt_Create]"
	var err error
	now := time.Now()
	if s, ok := claims.(jwt.StandardClaims); ok {
		s.IssuedAt = now.Unix()
		if j.ExpiresAt > 0 {
			s.ExpiresAt = j.ExpiresAt
		}
		if j.Issuer != "" {
			s.Issuer = j.Issuer
		}
		claims = s
	}
	if s, ok := claims.(jwt.MapClaims); ok {
		s["iat"] = now.Unix()
		if j.ExpiresAt > 0 {
			s["exp"] = j.ExpiresAt
		}
		if j.Issuer != "" {
			s["iss"] = j.Issuer
		}
		claims = s
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(j.SignKey))
	if err != nil {
		logx.Ex(ctx, tag, "create SignedString failed err:%+v, signKey:%s, claims:%+v",
			err, j.SignKey, claims)
		return "", errorx.Wrap500Response(err, errorx.JwtCreateSignStringErr, "")
	}
	return token, nil
}

func (j *jwter) Parse(ctx context.Context, tokenStr string, claims jwt.Claims) (*jwt.Token, error) {
	tag := "[jwtx_jwt_Parse]"
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logx.Ex(ctx, tag, "unexpected signing method token:%+v", token)
			return nil, errorx.ErrJwtSignMethod
		}
		return []byte(j.SignKey), nil
	})
	if err != nil {
		logx.Ex(ctx, tag, "jwt ParseWithClaims failed err:%+v, token:%+v", err, token)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return token, errorx.ErrJwtTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return token, errorx.ErrJwtTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return token, errorx.ErrJwtTokenNotValidYet
			} else {
				return token, errorx.ErrJwtTokenInvalid
			}
		}
	}
	return token, nil
}
