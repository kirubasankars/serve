package serve

import (
	"net/http"
	"strings"
)

type authenticator struct{}

func (auth *authenticator) Validate(ctx *Context, request *http.Request) bool {

	namespace := ctx.Namespace
	application := ctx.Application
	module := ctx.Module

	if namespace == nil || module == nil {
		return false
	}

	var roles map[string][]string
	if application != nil && application.Config != nil {
		roles = application.Config.Roles
	}
	if namespace.Config != nil {
		if roles == nil {
			roles = make(map[string][]string)
		}
		nroles := namespace.Config.Roles
		for k, v := range nroles {
			roles[k] = v
		}
	}

	if roles != nil {

		for _, v := range ctx.User.Roles {
			if role, done := roles[v]; done == true {
				for _, auth := range role {
					parts := strings.Split(auth, ":")
					if len(parts) == 2 {
						permissionName := parts[1]
						if module.Name == parts[0] {
							if exp, p := module.permittedRoutes[permissionName]; p == true && exp.MatchString(request.Method+" "+ctx.Path) {
								return true
							}
						}
					}
				}
			}
		}
		ctx.Module = nil
	}
	return false
}
