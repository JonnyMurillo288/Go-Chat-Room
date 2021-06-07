package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// template struct to hold out templates
// [templateName] *template.Template
type tpl struct {
	Templates map[string]*template.Template
}


// create a template struct by passing the name of the template
func newTemplate() *tpl {
	return &tpl{
		Templates: make(map[string]*template.Template),
	}
}

// add Template to the map in tpl struct
func (t *tpl) addTemplate(name string, temp *template.Template) {
	t.Templates[name] = temp
}

// create a template and automatically adds to tpl obj
// returns as well if needed
func (t *tpl) createTemplate(name string) *template.Template {
	t.Templates[name] = template.Must(template.New(name+".html").ParseFiles("templates/"+name+".html"))
	t.addTemplate(name,t.Templates[name])
	return t.Templates[name]
}

// handles new client tempalate to execute
// new client be POST METHOD for the client
func newClientHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Handling template")
		if err := t.Execute(w,nil); err != nil {
			http.Error(w, fmt.Sprintf("error executing template %s",err), http.StatusInternalServerError)
		}
	})
}

// function to save username and password in db
func saveLogin(){
	return
}


