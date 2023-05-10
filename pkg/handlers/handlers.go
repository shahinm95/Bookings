package handlers

import (
	"net/http"

	"github.com/shahinm95/bookings/pkg/config"
	"github.com/shahinm95/bookings/pkg/models"
	"github.com/shahinm95/bookings/pkg/render"
)

// for creating go temlpate file instead html => home.page.tmpl => create .tmpl file

// handlers should have access to appConfig because handlers should have access to all kind of configurations
// even though handlers may not use TemplateCache but give access to config file for future usecases
// like when we connect to database , we share database connection with config and using repository
// for this we use repository pattern : a commen pattern that allows us to swap compnents out of our application
// with minimal changes required to codebase

// repository used by handlers
var Repo *Repository

type Repository struct {
	App *config.AppCongif
}
 
// creates a new Repository
func NewRepo (a *config.AppCongif) *Repository {
	return &Repository{
		App :a,
	}
}

// it sets repository for handlers
func NewHandlers (r *Repository) {
	Repo = r
}
// we these function in main file above to get access to Appconfig then store in Repo 
// then we give stored data in Repo to handler via recievers



//defining handlers
// by giving a reciever to handler we giver access to AppConfig data
func (m *Repository) Home (w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.SessionManager.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
	
}

func (m *Repository) About (w http.ResponseWriter, r *http.Request){
	//preform some logic
	stringMap := map[string]string{}
	stringMap["test"]= "This is from handlers"
	// send data to template

	remoteIp := m.App.SessionManager.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
	
}
