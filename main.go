package main

import (
	"github.com/gorilla/mux"
	"net/http"
	handler "./handlers"
)
var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", handler.IndexPageHandler)
	router.HandleFunc("/checkuser", handler.CheckUserHandler)
	router.HandleFunc("/login", handler.LoginPageHandler)
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}