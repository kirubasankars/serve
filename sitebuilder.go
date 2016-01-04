package main

type SiteBuilder interface {
	Build(site *Site)
}

type CommonSiteBuilder struct {
}

func (builder *CommonSiteBuilder) Build(site *Site) {
	site.SetupHandler(site.uri, new(FileHandler))
	site.SetupHandler(site.uri+"/", new(FileHandler))
	site.SetupHandler(site.uri+"/api/", new(ApiHandler))
	site.SetupHandler(site.uri+"/html/", new(HtmlTemplateHandler))
	site.SetupHandler(site.uri+"/text/", new(TextTemplateHandler))
}

type AuthSiteBuilder struct {
}

func (builder *AuthSiteBuilder) Build(site *Site) {
	site.SetupHandler(site.uri, new(AuthSiteHandler))
	site.SetupHandler(site.uri+"/", new(AuthSiteHandler))
}
