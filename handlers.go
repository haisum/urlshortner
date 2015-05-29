// This is simple web app for url shortening. Code available at: http://github.com/haisum/urlshortner
//
// Features
//
// 	- Integrated db drivers and http  server, so you just run the binary after compilation without need of installing anything else.
// 	- Give a long url to get shortened link
// 	- Re-captcha for non registered users
// 	- Sqlite3 database support
// 	- Ability to register and login so you could come back to your list of shortened urls
// 	- Stats about clicks, geography and history of url etc

package urlshortner

import (
	"html/template"
	"log"
	"net/http"
)

// Errors type is used to record all errors ocurred during different operations in this web app and is printed as json to browser.
type Errors struct {
	errors []error
}

// Handles all GET requests intended for "/" route
// Renders form for shortening url
// If user is logged in, it also renders all previosuly shortened links of user with some stats
func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatalf("Could not find template file template/index.html. Error: %s.", err)
	}
	t.Execute(rw, nil)
}

// This handler receives POST request with parameter "url" on route /shorten
// and prints json for Url struct object or Errors struct object for errors
func ShortenHandler(rw http.ResponseWriter, req *http.Request) {

}

// This handler receives GET request with parameter in format /l/{id}
// and redirects to longer version of {id} if url is shortened or returns a 404 if {id} is not found
func ExpandHandler(rw http.ResponseWriter, req *http.Request) {

}

// This renders stats template and supplies it with stats for a particular shortened url.
// This handler receives GET requests in format /stats/{id:[0-9]+}
func StatsHandler(rw http.ResponseWriter, req *http.Request) {

}

// Receives "email", "password" in a POST request, if match found, prints User struct otherwise prints Errors struct object.
// This also sets session via gorilla/sessions and stores User struct object in session "user"
// Format for request is /login
func LoginHandler(rw http.ResponseWriter, req *http.Request) {

}

// Receives POST request on route /logout , unsets session "user" and redirects to "/"
func LogoutHandler(rw http.ResponseWriter, req *http.Request) {

}

// Receives POST request on route /register, with form values "email" and "password".
// If email is valid, isn't already registered and password is at least 6 chars long handler inserts data in users table and
// calls LoginHandler. Otherwise, Errors struct object is printed as json.
func RegisterHandler(rw http.ResponseWriter, req *http.Request) {

}
