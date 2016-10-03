package serve

// User user object
type User struct {
	ID         string
	Password   string
	Namespaces map[string]*UserNamespace
}

//UserNamespace user namespace
type UserNamespace struct {
	Apps map[string]*UserApplication
}

//UserApplication user application
type UserApplication struct {
	Roles []string
}

//NewUser - will create user object.
func NewUser(id string) *User {
	user := new(User)
	return user
}
