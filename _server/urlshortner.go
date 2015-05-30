package main

import (
	"github.com/gorilla/mux"
	"github.com/haisum/urlshortner"
	"log"
	"net/http"
)

func main() {
	db := urlshortner.Db{
		Name: "./urlshortner.db",
	}
	db.ConnectDb()
	r := mux.NewRouter()
	r.HandleFunc("/", urlshortner.HomeHandler).Methods("GET")
	r.HandleFunc("/stats/{id:[0-9]+}", urlshortner.StatsHandler).Methods("GET")
	r.HandleFunc("/l/{url}", urlshortner.ExpandHandler).Methods("GET")
	r.HandleFunc("/shorten", urlshortner.ShortenHandler).Methods("POST")
	r.HandleFunc("/login", urlshortner.LoginHandler).Methods("POST")
	//yes, logout should be a post request! See: http://blog.codinghorror.com/cross-site-request-forgeries-and-you/
	r.HandleFunc("/logout", urlshortner.LoginHandler).Methods("POST")
	r.HandleFunc("/register", urlshortner.RegisterHandler).Methods("POST")
	//serve static content from static folder
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	http.Handle("/", r)

	log.Print("\nStarting server on port 8081. You can browse web app on http://localhost:8081\n")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Error ocurred: %s", err)
	}
}
