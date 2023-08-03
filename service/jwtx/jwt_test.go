package jwtx

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTCustomize(t *testing.T) {
	type AuthJWTClaims struct {
		jwt.StandardClaims
		UserID     uint64
		Authorized bool
	}
	ctx := context.Background()
	j := NewJwt("123456")
	claims := AuthJWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
		uint64(1),
		true,
	}
	tokenStr, err := j.Create(ctx, claims)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenStr)
	claims2 := AuthJWTClaims{}
	token, err := j.Parse(ctx, tokenStr, &claims2)
	assert.Nil(t, err)
	assert.True(t, token.Valid)
	t.Log(claims2)
	t.Log(token.Claims)
}

func TestJWTMapClaims(t *testing.T) {
	ctx := context.Background()
	claims := jwt.MapClaims{}
	claims["user_id"] = uint64(1)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	j := NewJwt("123456")
	tokenStr, err := j.Create(ctx, claims)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenStr)
	t.Logf("token: %s", tokenStr)
	claims2 := jwt.MapClaims{}
	token, err := j.Parse(ctx, tokenStr, &claims2)
	assert.Nil(t, err)
	if c, ok := token.Claims.(*jwt.MapClaims); !ok {
		t.Errorf("token.Claims type must be %s, got %+v", "jwtgo.MapClaims", reflect.TypeOf(token.Claims))
	} else {
		assert.Nil(t, c.Valid())
		assert.Equal(t, claims["user_id"], uint64((*c)["user_id"].(float64)))
		assert.Equal(t, claims["authorized"], (*c)["authorized"])
	}
	assert.Equal(t, claims["user_id"].(uint64), uint64(claims2["user_id"].(float64)))
	assert.Equal(t, claims["authorized"], claims2["authorized"])
	t.Log(token.Claims)
	t.Log(claims2)
}
