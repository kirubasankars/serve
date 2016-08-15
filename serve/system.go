package serve

//System inferface for driver
type System interface {
	Build(ctx *Context, uri string)
}
