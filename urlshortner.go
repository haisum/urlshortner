package main

import (
	"github.com/gorilla/mux"
	"github.com/haisum/urlshortner"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", urlshortner.HomeHandler).Methods("GET")
	r.HandleFunc("/static/", http.FileServer(http.Dir("static"))).Methods("GET")
	r.HandleFunc("/stats/{id:[0-9]+}", urlshortner.StatsHandler).Methods("GET")
	r.HandleFunc("/l/{url}", urlshortner.ExpandHandler).Methods("GET")
	r.HandleFunc("/shorten", urlshortner.ShortenHandler).Methods("POST")
	r.HandleFunc("/login", urlshortner.LoginHandler).Methods("POST")
	//yes, logout should be a post request! See: http://blog.codinghorror.com/cross-site-request-forgeries-and-you/
	r.HandleFunc("/logout", urlshortner.LoginHandler).Methods("POST")
	r.HandleFunc("/register", urlshortner.RegisterHandler).Methods("POST")
	http.Handle("/", r)
}
