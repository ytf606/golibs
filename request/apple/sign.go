package apple

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"git.100tal.com/wangxiao_monkey_tech/lib/errorx"
	"git.100tal.com/wangxiao_monkey_tech/lib/logx"
	"git.100tal.com/wangxiao_monkey_tech/lib/pkg/osx"
	"git.100tal.com/wangxiao_monkey_tech/lib/pkg/stringx"
	"git.100tal.com/wangxiao_monkey_tech/lib/request"
	"github.com/golang-jwt/jwt/v4"
)

const (
	// ValidationURL is the endpoint for verifying tokens
	ValidationURL string = "https://appleid.apple.com/auth/token"
	// RevokeURL is the endpoint for revoking tokens
	RevokeURL string = "https://appleid.apple.com/auth/revoke"
	// ContentType is the one expected by Apple
	ContentType string = "application/x-www-form-urlencoded"
	// UserAgent is required by Apple or the request will fail
	UserAgent string = "go-lib-with-apple"
	// AcceptHeader is the content that we are willing to accept
	AcceptHeader string = "application/json"

	//jwt Audience
	AudienceURL string = "https://appleid.apple.com"
	//jwt Algorithm
	ES256Alg string = "ES256"
)

// 发送Header信息
var Header = map[string]string{
	"Content-Type": ContentType,
	"accept":       AcceptHeader,
	"User-Agent":   UserAgent,
}

// SignClient is an interface to call the validation API
type SignClient interface {
	VerifyWebToken(ctx context.Context, code, redirectUri string) (result ValidationResponse, err error)
	VerifyAppToken(ctx context.Context, code string) (result ValidationResponse, err error)
	VerifyRefreshToken(ctx context.Context, refreshToken string) (result RefreshResponse, err error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) (result RevokeResponse, err error)
	RevokeAccessToken(ctx context.Context, accessToken string) (result RevokeResponse, err error)
}

// Client implements ValidationClient
type Client struct {
	validationURL string
	revokeURL     string
	SignKey       string
	TeamId        string
	ClientId      string
	KeyId         string
	Secret        string
}

func NewAppleSign(signKey, teamId, clientId, keyId string) *Client {
	client := &Client{
		validationURL: ValidationURL,
		revokeURL:     RevokeURL,
		SignKey:       signKey,
		TeamId:        teamId,
		ClientId:      clientId,
		KeyId:         keyId,
	}
	return client
}

/*
GenerateClientSecret generates the client secret used to make requests to the validation server.
The secret expires after 6 months

signingKey - Private key from Apple obtained by going to the keys section of the developer section
teamID - Your 10-character Team ID
clientID - Your Services ID, e.g. com.aaronparecki.services
keyID - Find the 10-char Key ID value from the portal
*/
func (c *Client) GenSecret() error {
	if c.Secret != "" {
		return nil
	}
	block, _ := pem.Decode([]byte(c.SignKey))
	if block == nil {
		return errorx.New500Response(errorx.PemParseCodeErr, "parse result empty after pem decode")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	// Create the Claims
	now := time.Now()
	claims := &jwt.StandardClaims{
		Issuer:    c.TeamId,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Hour*24*180 - time.Second).Unix(), // 180 days
		Audience:  AudienceURL,
		Subject:   c.ClientId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["alg"] = ES256Alg
	token.Header["kid"] = c.KeyId

	if c.Secret, err = token.SignedString(privKey); err != nil {
		return err
	}
	return nil
}

// VerifyWebToken sends the WebValidationTokenRequest and gets validation result
// Code is the authorization code received from your application’s user agent.
// The code is single use only and valid for five minutes.
// RedirectURI is the destination URI the code was originally sent to.
// Redirect URLs must be registered with Apple. You can register up to 10. Apple will throw an error with IP address
// URLs on the authorization screen, and will not let you add localhost in the developer portal.
func (c *Client) VerifyWebToken(ctx context.Context, code, redirectUri string) (result ValidationResponse, err error) {
	if err = c.GenSecret(); err != nil {
		return
	}

	body := map[string]string{
		"client_id":     c.ClientId,
		"client_secret": c.Secret,
		"code":          code,
		"redirect_uri":  redirectUri,
		"grant_type":    "authorization_code",
	}

	resp, err := request.New().PostForm(ctx, c.validationURL, body, Header)
	if err != nil {
		return
	}
	if err = stringx.Decoder(resp, &result); err != nil {
		logx.Ex(ctx, osx.PF(), "json decode apple verify web token response failed err:%+v, rawData:%+v", err, string(resp))
		return
	}
	if result.Error != "" {
		logx.Ex(ctx, osx.PF(), "apple verify web token response raw result:%+v", result)
		err = errorx.New500Response(errorx.AppleTokenInvalidErr, fmt.Sprintf("apple verify web token response raw info:%+v", result))
	}
	return
}

// VerifyAppToken sends the AppValidationTokenRequest and gets validation result
// The authorization code received in an authorization response sent to your app. The code is single-use only and valid for five minutes.
// Authorization code validation requests require this parameter.
func (c *Client) VerifyAppToken(ctx context.Context, code string) (result ValidationResponse, err error) {
	if err = c.GenSecret(); err != nil {
		return
	}

	body := map[string]string{
		"client_id":     c.ClientId,
		"client_secret": c.Secret,
		"code":          code,
		"grant_type":    "authorization_code",
	}

	resp, err := request.New().PostForm(ctx, c.validationURL, body, Header)
	if err != nil {
		return
	}
	if err = stringx.Decoder(resp, &result); err != nil {
		logx.Ex(ctx, osx.PF(), "json decode apple verify app token response failed err:%+v, rawData:%+v", err, string(resp))
		return
	}
	if result.Error != "" {
		logx.Ex(ctx, osx.PF(), "apple verify app token response raw result:%+v", result)
		err = errorx.New500Response(errorx.AppleTokenInvalidErr, fmt.Sprintf("apple verify app token response raw info:%+v", result))
	}
	return
}

// VerifyRefreshToken sends the WebValidationTokenRequest and gets validation result
// RefreshToken is the refresh token given during a previous validation
func (c *Client) VerifyRefreshToken(ctx context.Context, refreshToken string) (result RefreshResponse, err error) {
	if err = c.GenSecret(); err != nil {
		return
	}

	body := map[string]string{
		"client_id":     c.ClientId,
		"client_secret": c.Secret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	resp, err := request.New().PostForm(ctx, c.validationURL, body, Header)
	if err != nil {
		return
	}
	if err = stringx.Decoder(resp, &result); err != nil {
		logx.Ex(ctx, osx.PF(), "json decode apple verify refresh token response failed err:%+v, rawData:%+v", err, string(resp))
		return
	}
	if result.Error != "" {
		logx.Ex(ctx, osx.PF(), "apple verify refresh token response raw result:%+v", result)
		err = errorx.New500Response(errorx.AppleTokenInvalidErr, fmt.Sprintf("apple verify refresh token response raw info:%+v", result))
	}
	return
}

// RevokeRefreshToken revokes the Refresh Token and gets the revoke result
// RefreshToken is the refresh token given during a previous validation
func (c *Client) RevokeRefreshToken(ctx context.Context, refreshToken string) (result RevokeResponse, err error) {
	if err = c.GenSecret(); err != nil {
		return
	}

	body := map[string]string{
		"client_id":       c.ClientId,
		"client_secret":   c.Secret,
		"token":           refreshToken,
		"token_type_hint": "refresh_token",
	}

	resp, err := request.New().PostForm(ctx, c.revokeURL, body, Header)
	if err != nil {
		return
	}
	if err = stringx.Decoder(resp, &result); err != nil {
		logx.Ex(ctx, osx.PF(), "json decode apple revoke refresh token response failed err:%+v, rawData:%+v", err, string(resp))
		return
	}
	if result.Error != "" {
		logx.Ex(ctx, osx.PF(), "apple revoke refresh token response raw result:%+v", result)
		err = errorx.New500Response(errorx.AppleTokenInvalidErr, fmt.Sprintf("apple revoke refresh token response raw info:%+v", result))
		return
	}
	return
}

// RevokeAccessToken revokes the Access Token and gets the revoke result
// AccessToken is the auth token given during a previous validation
func (c *Client) RevokeAccessToken(ctx context.Context, accessToken string) (result RevokeResponse, err error) {
	if err = c.GenSecret(); err != nil {
		return
	}

	body := map[string]string{
		"client_id":       c.ClientId,
		"client_secret":   c.Secret,
		"token":           accessToken,
		"token_type_hint": "access_token",
	}

	resp, err := request.New().PostForm(ctx, c.revokeURL, body, Header)
	if err != nil {
		return
	}
	if err = stringx.Decoder(resp, &result); err != nil {
		logx.Ex(ctx, osx.PF(), "json decode apple revoke access token response failed err:%+v, rawData:%+v", err, string(resp))
		return
	}
	if result.Error != "" {
		logx.Ex(ctx, osx.PF(), "apple revoke access token response raw result:%+v", result)
		err = errorx.New500Response(errorx.AppleTokenInvalidErr, fmt.Sprintf("apple revoke access token response raw info:%+v", result))
		return
	}
	return
}

// GetUniqueID decodes the id_token response and returns the unique subject ID to identify the user
func GetUniqueId(idToken string) (string, error) {
	j, err := GetClaims(idToken)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", j["sub"]), nil
}

// GetClaims decodes the id_token response and returns the JWT claims to identify the user
func GetClaims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errorx.ErrJwtTokenMalformed
	}
	var resp map[string]interface{}
	if err := stringx.RawUrlDecodeAndUnmarshal(parts[1], &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
