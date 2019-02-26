package handlers

import (
	"fmt"
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/PuerkitoBio/goquery"
	"strings" 
	"html/template"
)

type Course struct{
	ISO 	string
	Name	string
	Kurs	string
}

type ViewData struct{
	Courses []Course
	user string
}

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUser(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func parseUrl(url string) ViewData{
	doc, err := goquery.NewDocument(url)
	_check(err)
	
	var (
		rows [][]string
		row []string
	)
	course := Course{}
	courses := ViewData{}
	doc.Find("#myTable").Each(func (index int, table *goquery.Selection){
		table.Find("tr").Each(func (indexTr int, tr *goquery.Selection){
			tr.Find("td").Each(func (indexTd int, td *goquery.Selection){
				text := strings.Trim(td.Text(), "\r\n\t\v\f")
				if (text != ""){
					row = append(row, text)
				}
			})
			rows = append(rows, row)
			row = nil
		})
	})
	for _, val := range rows {
		if (len(val) > 0){
			if (val[1] == "810" || val[1] == "840" || val[1] == "978") {
				course.ISO 	= val[1]
				course.Name = val[3]
				course.Kurs = val[4]

				courses.Courses = append(courses.Courses, course)
			}
		}
	}
	return courses
}


func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	user := getUser(request)
	if user != "" {
		fmt.Println("Hello, ", user)
	} else {
		http.Redirect(response, request, "/login", 302)
	}
	
	login := request.FormValue("login")
	
	courses := parseUrl("http://nbt.tj/tj/kurs/kurs.php")
	newCourses := courses
	newCourses.user = login

	tmpl, err := template.ParseFiles("static/index.html") //parse the html file homepage.html
	if err != nil { // if there is an error
  		fmt.Println("Error parsing file")
  	}
    err = tmpl.Execute(response, newCourses) 
    if err != nil { 
		fmt.Println("Error executing template")
  	}
}
