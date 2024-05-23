package main

import (
	"fmt"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p Page) save() error {
  // This create a copy of p, modify the copy then use it to write a file
	filename := p.Title + ".txt"
  p.Body = []byte("Body modified.")
	return os.WriteFile(filename, p.Body, 0600)
}

func (p *Page) savePtr() error {
  // This takes the p, modify it and write it as a file
	filename := p.Title + ".txt"
  p.Body = []byte("Body modified.")
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page , error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
  if err != nil {
    // return Page{Title: "FailPage", Body: []byte("There is an error")}, err
    return nil, err
  }
	return &Page{Title: title, Body: body}, nil
}

func main() {
  p1 := Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  fmt.Println(string(p1.Body))
  
  p1.save()
  fmt.Println("Using save()", string(p1.Body))

  p1.savePtr()
  fmt.Println("Using savePtr()", string(p1.Body))
  
  p2, err := loadPage("TestPage")
  fmt.Println(err)
  fmt.Println(string(p2.Body))

}
