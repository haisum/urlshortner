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
	"fmt"
	"github.com/gorilla/mux"
	"github.com/haisum/recaptcha"
	"github.com/haisum/urlshortner/timediff"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//template data is used to pass data to template index.html
type TemplateData struct {
	Email string
}

// Handles all GET requests intended for "/" route
// Renders form for shortening url
// If user is logged in, it also renders all previosuly shortened links of user with some stats
func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	//t, err := template.Delims("{#", "#}").ParseFiles("./templates/index.html")
	t := template.New("index.html")
	t.Delims("[[", "]]")
	t, err := t.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatalf("Could not find template file template/index.html. Error: %s.", err)
	}

	session := app.GetSession(req)
	td := TemplateData{}
	if u, ok := session.Values["user"].(User); ok {
		td.Email = u.Email
	}
	err = t.Execute(rw, td)
	if err != nil {
		log.Printf("Couldn't show template with data %v. Error %s", td, err)
	}
}

type ShortenResponse struct {
	Errors  []string
	Success bool
	Url     string
	LongUrl string
}

// This handler receives POST request with parameter "url" on route /shorten
// and prints json for Url struct object or Errors struct object for errors
func ShortenHandler(rw http.ResponseWriter, req *http.Request) {
	//just to show off the ajax loading effect
	//remove this in production
	time.Sleep(time.Second * 1)
	session := app.GetSession(req)
	ul := Url{
		Url: req.PostFormValue("url"),
	}
	r := ShortenResponse{
		Errors: make([]string, 0),
	}
	if u, ok := session.Values["user"].(User); ok {
		ul.Userid = u.Id
	} else {
		re := recaptcha.R{
			Secret: "6LdRlAcTAAAAAGIIlUc0_jlqgnOMsr0AaCTkS-hg",
		}
		valid := re.Verify(*req)
		if !valid {
			log.Printf("Re-captcha attempt failed. %v", re.LastError())
			r.Errors = append(r.Errors, "You didn't prove you're a human.")
			fmt.Fprintf(rw, "%s", GetJson(r))
			return
		}
	}

	errors := ul.Validate()

	if errors != nil {
		for _, err := range errors {
			r.Errors = append(r.Errors, err.Error())
		}
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	err := ul.Save()

	if err != nil {
		log.Printf("Couldn't save url %s", err)
		r.Errors = append(r.Errors, "A server error occurred. Contact admin for support.")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	r.Success = true
	r.Url = strings.TrimRight(req.Referer(), "/") + "/l/" + IdToUrlString(ul.Id)
	r.LongUrl = ul.Url
	fmt.Fprintf(rw, "%s", GetJson(r))
}

// This handler receives GET request with parameter in format /l/{url}
// and redirects to longer version of {url} if url is shortened or returns a 404 if {url} is not found
func ExpandHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if b36id, ok := vars["url"]; ok {
		id, err := UrlStringToId(b36id)
		if err != nil {
			log.Printf("Couldn't convert %s to id.", b36id)
			fmt.Printf("Url Not Found.\n")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		url := Url{
			Id: id,
		}
		err = url.Get()
		if err != nil {
			log.Printf("Url %s => %d not found.", b36id, url.Id)
			fmt.Printf("Url Not Found.\n")
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		hit := Hit{
			Ip:       req.RemoteAddr,
			Referrer: req.Referer(),
			Urlid:    url.Id,
		}
		err = hit.Save()
		if err != nil {
			log.Printf("Couldn't save hit %v. Error: %s", hit, err)
		}
		http.Redirect(rw, req, url.Url, http.StatusMovedPermanently)
	} else {
		log.Printf("b36id not supplied. %v", vars)
	}
}

type LoginResponse struct {
	Errors  []string
	Success bool
}

// Receives "email", "password" in a POST request, if match found, prints User struct otherwise prints Errors struct object.
// This also sets session via gorilla/sessions and stores User struct object in session "user"
// Format for request is /login
func LoginHandler(rw http.ResponseWriter, req *http.Request) {
	//just to show off the ajax loading effect
	//remove this in production
	time.Sleep(time.Second * 1)
	r := LoginResponse{
		Errors: make([]string, 0),
	}

	session := app.GetSession(req)
	if u, ok := session.Values["user"].(User); ok {
		log.Printf("Already logged in as %s.\n", u.Email)
		r.Errors = append(r.Errors, "Already logged in as "+u.Email)
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	email := req.PostFormValue("email")
	password := req.PostFormValue("password")

	if strings.TrimSpace(email) == "" || password == "" {
		r.Errors = append(r.Errors, "Email or password can't be blank.")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	if app.HasFailedAttempts(email) {
		log.Printf("Preventing %s from login. More than 5 failed attempts in last 30 secs.\n", email)
		r.Errors = append(r.Errors, "Too many failed attempts. Wait 30 seconds before retrying.")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	user := User{
		Email: email,
	}

	err := user.EmailGet()
	if err != nil {
		app.RecordFailedAttempt(email)
		log.Printf("Couldn't get email %s for ip %s. %s\n", email, req.RemoteAddr, err)
		r.Errors = append(r.Errors, "Wrong email or password.")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		app.RecordFailedAttempt(email)
		log.Printf("Password mismatch %s for ip %s. %s\n", email, req.RemoteAddr, err)
		r.Errors = append(r.Errors, "Wrong email or password.")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}
	r.Success = true
	//now save user in session
	user.Password = ""
	session.Values["user"] = user
	err = session.Save(req, rw)
	if err != nil {
		log.Printf("Error in creating session. %s", err)
	}

	fmt.Fprintf(rw, "%s", GetJson(r))
}

// Receives POST request on route /logout , unsets session "user" and redirects to "/"
func LogoutHandler(rw http.ResponseWriter, req *http.Request) {
	//just to show off the ajax loading effect
	//remove this in production
	time.Sleep(time.Second * 1)
	session := app.GetSession(req)
	delete(session.Values, "user")
	session.Save(req, rw)
	http.Redirect(rw, req, req.Referer(), http.StatusFound)
}

type RegisterResponse struct {
	Errors  []string
	Success bool
	User    User
}

// Receives POST request on route /register, with form values "email" and "password".
// If email is valid, isn't already registered and password is at least 6 chars long handler inserts data in users table and
// calls LoginHandler. Otherwise, Errors struct object is printed as json.
func RegisterHandler(rw http.ResponseWriter, req *http.Request) {
	//just to show off the ajax loading effect
	//remove this in production
	time.Sleep(time.Second * 1)
	r := RegisterResponse{
		Errors: make([]string, 0),
	}

	session := app.GetSession(req)
	if u, ok := session.Values["user"].(User); ok {
		log.Printf("Attempt to register! Already Logged in as %s.", u.Email)
		r.Errors = append(r.Errors, "You are logged in as "+u.Email+". Logout to continue.\n")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	email := req.PostFormValue("email")
	password := req.PostFormValue("password")

	if strings.TrimSpace(email) == "" || password == "" {
		r.Errors = append(r.Errors, "Email or password can't be blank.")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	user := User{
		Email:    email,
		Password: password,
	}

	errs := user.Validate()
	if errs != nil {
		for _, err := range errs {
			r.Errors = append(r.Errors, err.Error())
		}
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}

	user.Save()
	user.Password = ""
	session.Values["user"] = user
	err := session.Save(req, rw)
	if err != nil {
		log.Printf("Error in creating session. %s", err)
	}
	log.Printf("%v", session.Values)

	r.User = user
	r.Success = true
	fmt.Fprintf(rw, "%s", GetJson(r))
}

type UrlResponse struct {
	Hits    []Hit
	Created string
	Id      string
	Url     string
}

type MeResponse struct {
	Success bool
	Email   string
	Urls    []UrlResponse
	Total   int64
}

func MeHandler(rw http.ResponseWriter, req *http.Request) {
	session := app.GetSession(req)
	r := MeResponse{
		Urls: make([]UrlResponse, 0),
	}
	user, ok := session.Values["user"].(User)

	if !ok {
		log.Printf("User not registered. Attempt to access /me")
		fmt.Fprintf(rw, "%s", GetJson(r))
		return
	}
	r.Email = user.Email
	r.Success = true

	offset, err := strconv.Atoi(req.FormValue("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = 10
	}

	urls, err := GetAllUrls(user.Id, offset, limit)

	if err != nil {
		log.Printf("Couldn't get urls. %s", err)
	}

	for _, u := range urls {
		hits, err := u.GetHits()
		if err != nil {
			log.Printf("Couldn't get hits. %s", err)
		}
		t := time.Unix(u.Ondate, 0)
		ur := UrlResponse{
			Url:     u.Url,
			Created: timediff.GetDifference(time.Now(), t),
			Id:      strings.TrimRight(req.Referer(), "/") + "/l/" + IdToUrlString(u.Id),
			Hits:    hits,
		}
		r.Urls = append(r.Urls, ur)
	}
	total, err := GetTotalUrls(user.Id)
	if err != nil {
		log.Printf("Couldn't get total urls of user %d. %s", user.Id, err)
	}
	r.Total = total
	fmt.Fprintf(rw, "%s", GetJson(r))
}
