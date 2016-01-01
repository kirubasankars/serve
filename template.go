package main

import (
	"net/http"
	"strings"
)

type HtmlServe struct {
	site *Site
}

func (htmlServe *HtmlServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var site = htmlServe.site
	parts := strings.Split(r.URL.Path, "/html")
	model := site.Model(parts[1] + "/" + strings.ToLower(r.Method))
	site.HandleHTMLTemplate(model, w, r)
}

type TextServe struct {
	site *Site
}

func (textServe *TextServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var site = textServe.site
	parts := strings.Split(r.URL.Path, "/text")
	model := site.Model(parts[1] + "/" + strings.ToLower(r.Method))
	site.HandleTextTemplate(model, w, r)
}
