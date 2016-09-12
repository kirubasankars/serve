package serve

//AccessToken storage
type AccessToken struct {
	Token  string
	Client *Client
	User   *User
}
