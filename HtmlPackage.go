package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
//   m := validPath.FindStringSubmatch(r.URL.Path)
//   if m == nil {
//     http.NotFound(w, r)
//     return "", errors.New("Invalid Page Title")
//   }
//   return m[2], nil
// }

type Page struct {
	Title string
	Body  []byte
}

func (p Page) save(body string) error {
	filename := p.Title + ".txt"
	p.Body = []byte(body)
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Println("ERROR ON EXECUTE", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Println(m)
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	//   return
	// }
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	//   return
	// }
	p, err := loadPage(title)
	if err != nil {
		fmt.Println("ERROR ON HANDLER")
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// title, err := getTitle(w, r)
	// if err != nil {
	//   return
	// }
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save(body)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	// http.Handlemlnc("/view/", viewHandler)
	// http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	fmt.Println("Listening to port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
