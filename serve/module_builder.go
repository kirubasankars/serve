package serve

// ModuleHandlerProvider interface
type ModuleHandlerProvider interface {
	Build(module Module)
}
