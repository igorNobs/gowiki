package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("./templates/edit.html", "./templates/view.html"))

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	_ = os.Mkdir("documents", 0777)
	filename := "./documents/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./documents/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body, Title: title}, nil
}

func renderTemplate(writer http.ResponseWriter, templateName string, p *Page) {
	err := templates.ExecuteTemplate(writer, templateName+".html", p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
