package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
	"html/template"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

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

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}
func _check(err error) {
	if err != nil {
		panic(err)
	}
}
func parseUrl(url string) [][]string{
	doc, err := goquery.NewDocument(url)
	_check(err)
	
	var (
		rows [][]string
		row []string
		courses [][]string
	)
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
				courses = append(courses, val)
			}
		}
	}
	return courses
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	user := getUser(request)
	if user != "" {
		fmt.Println("Hello, ", user)
	} else {
		http.Redirect(response, request, "/login", 302)
	}
	
	login := request.FormValue("login")
	pass := request.FormValue("pass")
	fmt.Println("login = ", login);
	fmt.Println("pass = ", pass);

	tmpl, err := template.ParseFiles("static/index.html") //parse the html file homepage.html
	if err != nil { // if there is an error
  		fmt.Println("Error parsing file")
  	}
    err = tmpl.Execute(response, nil) 
    if err != nil { 
		fmt.Println("Error executing template")
  	}
}


func loginPageHandler(response http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("static/login.html") //parse the html file homepage.html
	if err != nil { // if there is an error
  		fmt.Println("Error parsing file")
  	}
    err = tmpl.Execute(response, nil) 
    if err != nil { 
		fmt.Println("Error executing template")
	}
	//http.Redirect(response, request, "/checkuser", 302)

}

func checkUserHandler(response http.ResponseWriter, request *http.Request) {
	login := request.FormValue("login")
	pass := request.FormValue("pass")
	if login != "" && pass != "" {
		// .. check credentials ..
		setSession(login, response)
	}
	http.Redirect(response, request, "/", 302)
}

var router = mux.NewRouter()

func main() {
	parseUrl("http://nbt.tj/tj/kurs/kurs.php")
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/checkuser", checkUserHandler)
	router.HandleFunc("/login", loginPageHandler)
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}