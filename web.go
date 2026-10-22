package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
)

var coreNav = []NavItem{
	{Path: "/budget/", Label: "nav.budget"},
	{Path: "/stats", Label: "nav.stats"},
}

// parsePages pairs every template with layout.html, keyed by file name without extension.
func parsePages() (map[string]*template.Template, error) {
	funcs := template.FuncMap{"hasPrefix": strings.HasPrefix}
	files, err := fs.Glob(embedded, "templates/*.html")
	if err != nil {
		return nil, err
	}
	pages := map[string]*template.Template{}
	for _, f := range files {
		name := path.Base(f)
		if name == "layout.html" {
			continue
		}
		t, err := template.New(name).Funcs(funcs).ParseFS(embedded, "templates/layout.html", f)
		if err != nil {
			return nil, err
		}
		pages[strings.TrimSuffix(name, ".html")] = t
	}
	return pages, nil
}

func (a *App) render(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	data["L"] = a.locale(r)
	data["User"] = currentUser(r)
	data["Path"] = r.URL.Path
	data["Nav"] = a.nav()
	data["LoginEnabled"] = a.loginEnabled()
	var buf bytes.Buffer
	if err := a.pages[page].ExecuteTemplate(&buf, "layout", data); err != nil {
		a.fail(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func (a *App) nav() []NavItem {
	nav := append([]NavItem{}, coreNav...)
	for _, p := range plugins {
		nav = append(nav, p.Nav()...)
	}
	return nav
}

func (a *App) fail(w http.ResponseWriter, err error) {
	log.Print(err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}
