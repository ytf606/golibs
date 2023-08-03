package apple

// ValidationResponse is based off of https://developer.apple.com/documentation/signinwithapplerestapi/tokenresponse
type ValidationResponse struct {
	// (Reserved for future use) A token used to access allowed data. Currently, no data set has been defined for access.
	AccessToken string `json:"access_token"`

	// The type of access token. It will always be "bearer".
	TokenType string `json:"token_type"`

	// The amount of time, in seconds, before the access token expires. You can revalidate with the "RefreshToken"
	ExpiresIn int `json:"expires_in"`

	// The refresh token used to regenerate new access tokens. Store this token securely on your server.
	// The refresh token isn’t returned when validating an existing refresh token. Please refer to RefreshReponse below
	RefreshToken string `json:"refresh_token"`

	// A JSON Web Token that contains the user’s identity information.
	IdToken string `json:"id_token"`

	// Used to capture any error returned by the endpoint. Do not trust the response if this error is not nil
	Error string `json:"error"`

	// A more detailed precision about the current error.
	ErrorDescription string `json:"error_description"`
}

// RefreshResponse is a subset of ValidationResponse returned by Apple
type RefreshResponse struct {
	// (Reserved for future use) A token used to access allowed data. Currently, no data set has been defined for access.
	AccessToken string `json:"access_token"`

	// The type of access token. It will always be "bearer".
	TokenType string `json:"token_type"`

	// The amount of time, in seconds, before the access token expires. You can revalidate with this token
	ExpiresIn int `json:"expires_in"`

	// Used to capture any error returned by the endpoint. Do not trust the response if this error is not nil
	Error string `json:"error"`

	// A more detailed precision about the current error.
	ErrorDescription string `json:"error_description"`
}

// RevokeResponse is based of https://developer.apple.com/documentation/sign_in_with_apple/revoke_tokens
type RevokeResponse struct {
	// Used to capture any error returned by the endpoint
	Error string `json:"error"`

	// A more detailed precision about the current error.
	ErrorDescription string `json:"error_description"`
}
