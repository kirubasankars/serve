package serve

// WebserverFlow WebserverFlow
type WebserverFlow interface {
	GetAuthorizationCode(clientID string, redirectURL string, username string, password string) string
	GetAccessToken(authorizationCode string, clientID string, clientSecret string, redirectURL string) *AccessToken
}

// UserAgentFlow UserAgentFlow
type UserAgentFlow interface {
	GetAccessToken(clientID string, redirectURL string, username string, password string) *AccessToken
}

// UserPasswordFlow UserPasswordFlow
type UserPasswordFlow interface {
	GetAccessToken(clientID string, clientSecret string, username string, password string) *AccessToken
}

// RefreshTokenFlow RefreshTokenFlow
type RefreshTokenFlow interface {
	GetAccessToken(refreshToken string, clientID string, clientSecret string) *AccessToken
}

//OAuth2Provider oauth
type OAuth2Provider interface {
	Init(server *Server)
	GetWebserver() WebserverFlow
	GetUserAgent() UserAgentFlow
	GetUserPassword() UserPasswordFlow
	GetRefreshToken() RefreshTokenFlow
}
