package serve

// Authentication forms authentication
type Authentication struct{}

// Authenticate authenticate with username and password
func (auth *Authentication) Authenticate(param AuthParameter) int {
	var count int

	if param.refreshToken == "refresh_token" {
		count++
	}
	if param.authorizationCode == "code" {
		count++
	}
	if param.username == "admin" && param.password == "admin" {
		count++
	}
	if param.clientID == "client_id" && param.redirectURI == "url" {
		count++
	}
	if param.clientID == "client_id" && param.clientSecret == "client_secret" {
		count++
	}

	return count
}

// AuthParameter authParameter
type AuthParameter struct {
	refreshToken      string
	authorizationCode string

	clientID     string
	redirectURI  string
	clientSecret string

	username string
	password string
}
