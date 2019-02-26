package handlers

import (
	"html/template"
	"fmt"
	"net/http"
)


func LoginPageHandler(response http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("static/login.html") //parse the html file homepage.html
	if err != nil { // if there is an error
  		fmt.Println("Error parsing file")
  	}
    err = tmpl.Execute(response, nil) 
    if err != nil { 
		fmt.Println("Error executing template")
	}
}
