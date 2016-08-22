package serve

// NamespaceConfigration config for namespace object
type NamespaceConfigration struct{}

// ApplicationConfigration config for application object
type ApplicationConfigration struct {
	Modules []string
	Roles   map[string][]string
}

// ModuleConfigration config for module object
type ModuleConfigration struct {
	Permissions map[string][]string
}
