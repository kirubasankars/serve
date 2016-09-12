package serve

//System inferface for driver
type System interface {
	GetConfig(path string) *[]byte

	Build(ctx *Context, uri string)

	GetOAuthProvider(server *Server) OAuth2Provider
}
