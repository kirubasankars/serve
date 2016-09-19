package driver

import (
	"encoding/json"
	"path/filepath"

	"github.com/kirubasankars/serve/serve"
)

//Webserver driver
type Webserver struct {
	server *serve.Server
}

// GetAuthorizationCode get authorization code
func (ws *Webserver) GetAuthorizationCode(clientID string, redirectURL string, username string, password string) string {
	return "12345678"
}

// GetAccessToken get access token
func (ws *Webserver) GetAccessToken(authorizationCode string, clientID string, clientSecret string, redirectURL string) *serve.AccessToken {
	res := new(serve.AccessToken)
	return res
}

//UserAgent user agent flow
type UserAgent struct {
	server *serve.Server
}

// GetAccessToken get access token
func (ua *UserAgent) GetAccessToken(clientID string, redirectURL string, username string, password string) *serve.AccessToken {
	server := ua.server
	path := server.Path()
	system := server.System

	user := getUser(system, path, username)
	client := getClient(system, path, clientID)

	//ctx := new(serve.Context)

	//fmt.Println(redirectURL)

	//system.Build(ctx, redirectURL)

	//namespace, _ := user.Namespaces[ctx.Namespace.Name]
	//app, _ := namespace.Apps[ctx.Application.Name]

	//fmt.Println(namespace, app)

	//&& namespace != nil && app != nil

	var t *serve.AccessToken

	if user.Password == password && client != nil {
		token := "access_token"

		t = new(serve.AccessToken)
		t.Token = token
		t.Client = client
		t.User = user
	}

	return t
}

//UserPassword user password flow
type UserPassword struct {
	server *serve.Server
}

// GetAccessToken get access token
func (up *UserPassword) GetAccessToken(clientID string, clientSecret string, username string, password string) *serve.AccessToken {
	server := up.server
	path := server.Path()
	system := server.System

	var t *serve.AccessToken
	user := getUser(system, path, username)
	client := getClient(system, path, clientID)

	if client.Secret == clientSecret && user.Password == password {
		token := "access_token"

		t = new(serve.AccessToken)
		t.Token = token
		t.Client = client
		t.User = user
	}

	return t
}

//RefreshToken refresh token
type RefreshToken struct {
	server *serve.Server
}

// GetAccessToken get access token
func (rt *RefreshToken) GetAccessToken(refreshToken string, clientID string, clientSecret string) *serve.AccessToken {
	t := new(serve.AccessToken)
	return t
}

func getUser(system serve.System, path string, userID string) *serve.User {
	usersConfig := system.GetConfig(filepath.Join(path, "users.json"))
	var users map[string]serve.User
	if err := json.Unmarshal(*usersConfig, &users); err != nil {
		panic(err)
	}
	user, _ := users[userID]
	return &user
}

func getClient(system serve.System, path string, clientID string) *serve.Client {
	clientsConfig := system.GetConfig(filepath.Join(path, "clients.json"))
	var clients map[string]serve.Client
	if err := json.Unmarshal(*clientsConfig, &clients); err != nil {
		panic(err)
	}
	client, _ := clients[clientID]
	return &client
}

//Authentication token
type Authentication struct {
	server *serve.Server
}

// Validate get access token
func (auth *Authentication) Validate(accessToken string, clientID string) *map[string]string {
	return nil
}

type fileOAuth2Provider struct {
	server *serve.Server
}

func (fp *fileOAuth2Provider) GetWebserver() serve.WebserverFlow {
	o := new(Webserver)
	o.server = fp.server
	return o
}

func (fp *fileOAuth2Provider) GetUserAgent() serve.UserAgentFlow {
	o := new(UserAgent)
	o.server = fp.server
	return o
}

func (fp *fileOAuth2Provider) GetUserPassword() serve.UserPasswordFlow {
	o := new(UserPassword)
	o.server = fp.server
	return o
}

func (fp *fileOAuth2Provider) GetRefreshToken() serve.RefreshTokenFlow {
	o := new(RefreshToken)
	o.server = fp.server
	return o
}

func (fp *fileOAuth2Provider) Init(server *serve.Server) {
	fp.server = server
}

//GetOAuthProvider GetOAuthProvider
func (fs *FileSystem) GetOAuthProvider(server *serve.Server) serve.OAuth2Provider {
	oauthProvider := new(fileOAuth2Provider)
	oauthProvider.Init(server)
	return oauthProvider
}
