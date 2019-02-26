package handlers

import (
	"net/http"
)

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
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

func CheckUserHandler(response http.ResponseWriter, request *http.Request) {
	login := request.FormValue("login")
	pass := request.FormValue("pass")
	if login != "" && pass != "" {
		// .. check credentials ..
		setSession(login, response)
	}
	http.Redirect(response, request, "/", 302)
}