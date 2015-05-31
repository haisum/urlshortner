package main

import (
	"github.com/gorilla/mux"
	"github.com/haisum/urlshortner"
	"log"
	"net/http"
)

func main() {

	urlshortner.Start()

	//routing logic
	r := mux.NewRouter()
	r.HandleFunc("/", urlshortner.HomeHandler).Methods("GET")
	r.HandleFunc("/l/{url}", urlshortner.ExpandHandler).Methods("GET")
	r.HandleFunc("/shorten", urlshortner.ShortenHandler).Methods("POST")
	r.HandleFunc("/login", urlshortner.LoginHandler).Methods("POST")
	//yes, logout should be a post request! See: http://blog.codinghorror.com/cross-site-request-forgeries-and-you/
	r.HandleFunc("/logout", urlshortner.LogoutHandler).Methods("POST")
	r.HandleFunc("/register", urlshortner.RegisterHandler).Methods("POST")
	//handler for showing data for currently logged in user
	r.HandleFunc("/me", urlshortner.MeHandler).Methods("GET")
	//serve static content from static folder
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	http.Handle("/", r)

	// start server
	log.Print("\nStarting server on port 8081. You can browse web app on http://localhost:8081\n")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Error ocurred: %s", err)
	}
}
