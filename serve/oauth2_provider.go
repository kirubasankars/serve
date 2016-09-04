package serve

// WebserverFlow WebserverFlow
type WebserverFlow interface {
	GetAuthorizationCode(clientID string, redirectURL string, username string, password string) string
	GetAccessToken(authorizationCode string, clientID string, clientSecret string, redirectURL string) *map[string]string
}

// UserAgentFlow UserAgentFlow
type UserAgentFlow interface {
	GetAccessToken(clientID string, redirectURL string) *map[string]string
}

// UserPasswordFlow UserPasswordFlow
type UserPasswordFlow interface {
	GetAccessToken(clientID string, clientSecret string, username string, password string) *map[string]string
}

// RefreshTokenFlow RefreshTokenFlow
type RefreshTokenFlow interface {
	GetAccessToken(refreshToken string, clientID string, clientSecret string) *map[string]string
}

//OAuth2Provider oauth
type OAuth2Provider interface {
	GetWebserver() WebserverFlow
	GetUserAgent() UserAgentFlow
	GetUserPassword() UserPasswordFlow
	GetRefreshToken() RefreshTokenFlow
}
