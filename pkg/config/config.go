package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

//appconfig holds application configuration
type AppCongif struct {
	UseCache bool
	TemplateCache map[string]*template.Template
	InProduction bool
	SessionManager *scs.SessionManager
}