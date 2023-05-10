package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/shahinm95/bookings/pkg/config"
	"github.com/shahinm95/bookings/pkg/models"
)

func RenderTemplate1(w http.ResponseWriter, templateName string) {
	// adding layout template right after template name
	parseTemplate, _ := template.ParseFiles("./templates/" + templateName, "./templates/base.layout.tmpl")
	err := parseTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template :", err)
		return
	}
}


// defing a way to cache templates strings instead of reading them from disk each time user makes a request 
// first way which is easier one =>

var templateCache = make(map[string] *template.Template)

func RenderTemplate2 (w http.ResponseWriter, templateName string ) {
	var template *template.Template
	var err error

	// checking if there is already a template with this name cached
	_, isThere := templateCache[templateName]
	if !isThere {
		// adding a new template to the cache
		err = CreateTemplateCache(templateName)
		if err !=nil {
			fmt.Println(err)
		}
		fmt.Println("adding template to cache first time")
	} else {
		// existing template
		fmt.Println("using existing template")
	}

	template = templateCache[templateName]
	err = template.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateTemplateCache ( templateName string) error {
	templates := []string {
		fmt.Sprintf("./templates/%s", templateName),
		"./templates/base.layout.tmpl",
	}
	//parse the templates
	templ , err1 := template.ParseFiles(templates...)
	if err1 != nil {
		return err1
	}
	templateCache[templateName] =templ
	return nil
}



// third way of chaching template which is bit more complex



// defining a function to have access to Appconfig
var App *config.AppCongif
func NewTemplates (a *config.AppCongif) {
	 App = a
} // wil call this function in main.go so from main.go you can have access to Appconfig here

// definig a function to give a default data to be avaiable to every page of site withoutmanually adding it
func defaultData (td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate (w http.ResponseWriter, temple string, templateData *models.TemplateData) {
	// // create a template chach
	// templateCache, err := CreateTemplateCacheThird()
	// if err!= nil {
	// 	//killing application
	// 	log.Fatal(err)
	// 	return
	// }

	// now instead of creating a new template cache, getting from AppConfig
	var templateCache map[string]*template.Template
	if App.UseCache {
		templateCache = App.TemplateCache
	} else {
		templateCache , _ = CreateTemplateCacheThird()
	}

	//get reqested template from chach
	templateName , ok := templateCache[temple]
	if !ok {
		log.Fatal("error loading template")
	}

	//arbitary way of totur , making a buffer using bytes.Buffer
	// A byte buffer is used to improve performance when writing a stream of data. 
	// buf will hold bytes, we will tryto execute the walues got from that map but rather than doing directly
	// we try to execute buffer directly, the reasong of doing this is for finer grainde error checking
	buff := new(bytes.Buffer)
	templateData = defaultData(templateData)
	err := templateName.Execute(buff, templateData) // to know where the error is comming from when executing
	if err != nil {
		log.Println(err) // this will tell error is comming from value stored in the map
	}


	//redner the template
	_, err = buff.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
	
}




func CreateTemplateCacheThird () (map[string]*template.Template, error){
	myCache := map[string]*template.Template{}

	//get all files that end with .page.temple from ./templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}


	//ranging through all files ending with .page.temple
	for _, page := range pages {
		// to extract name of file from path address
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// check if there are any templates layout for this template page
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil { return myCache, err}

		// if there are any templates layouts 
		//len => length
		if len(matches) >0 {
			// checking any templateSets , might require some file with files ending with .layout.temple
			// so we use ParseGlobe to parse them add thmem to templateSet
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {return myCache, err}
		}

		// adding templateSet to myCache
		myCache[name] = templateSet
	}



	return myCache, nil
}