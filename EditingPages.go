package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Page struct {
  Title string
  Body []byte
}

func loadPage(title string) (*Page, error){
  filename := title + ".txt"
  body, err := os.ReadFile(filename)
  if err != nil {
    fmt.Println("ERROR on os")
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: "ERROR"}
  }
  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/edit/"):] 
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  fmt.Fprintf(w, "<h1>Editing %s</h1>"+
    "<form action=\"/save/%s\" method=\"POST\">"+
    "<textarea name=\"body\">%s</textarea><br>"+
    "<input type=\"submit\" value=\"Save\">"+
    "</form>",
    p.Title, p.Title, p.Body) 
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  log.Fatal(http.ListenAndServe(":8000",nil))
}
