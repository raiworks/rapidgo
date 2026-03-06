package router

import "html/template"

// DefaultFuncMap returns the template function map with framework helpers.
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"route": Route,
	}
}
