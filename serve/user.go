package serve

// User user object
type User struct {
	Roles []string
}

func newUser(id string) *User {
	user := new(User)
	return user
}
