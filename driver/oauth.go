package driver

import "github.com/kirubasankars/serve/serve"

//WebServer driver
type WebServer struct{}

// GetAuthorizationCode get authorization code
func (ws *WebServer) GetAuthorizationCode(clientID string, redirectURL string, username string, password string) string {
	return "12345678"
}

// GetAccessToken get access token
func (ws *WebServer) GetAccessToken(authorizationCode string, clientID string, clientSecret string, redirectURL string) *map[string]string {
	res := map[string]string{
		"access_token":  "access_token",
		"refresh_token": "refresh_token",
		"issued_at":     "issued_at",
		"signature":     "signature",
	}
	return &res
}

//UserAgent user agent flow
type UserAgent struct{}

// GetAccessToken get access token
func (ua *UserAgent) GetAccessToken(clientID string, redirectURL string) *map[string]string {
	return nil
}

//UserPassword user password flow
type UserPassword struct{}

// GetAccessToken get access token
func (up *UserPassword) GetAccessToken(clientID string, clientSecret string, username string, password string) *map[string]string {
	res := map[string]string{
		"access_token": "access_token",
		"issued_at":    "issued_at",
		"signature":    "signature",
	}
	return &res
}

//RefreshToken refresh token
type RefreshToken struct{}

// GetAccessToken get access token
func (rt *RefreshToken) GetAccessToken(refreshToken string, clientID string, clientSecret string) *map[string]string {
	res := map[string]string{
		"access_token": "access_token",
		"issued_at":    "issued_at",
		"signature":    "signature",
	}
	return &res
}

//Authentication token
type Authentication struct{}

// Validate get access token
func (auth *Authentication) Validate(accessToken string, clientID string) *map[string]string {
	return nil
}

//GetOAuthProvider GetOAuthProvider
func (fs *FileSystem) GetOAuthProvider() serve.OAuth2Provider {
	return nil
}
