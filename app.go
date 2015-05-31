package urlshortner

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

//App struct holds variables that are globally required by our app
//Db holds database connection
//FailedAttempts keeps count of failed login attempts via an email and slows down login for current day unless correct credentials are given
//UrlsToday keeps count of urls a registered user shortened on current day, and doesn't allow more than 200 in one day.

type FailedAttempt struct {
	Last  int64
	Count int
}

type App struct {
	Db             *sqlx.DB
	FailedAttempts map[string]*FailedAttempt
	UrlsToday      map[int64]int
	SessionStore   *sessions.CookieStore
}

var app = new(App)

func Start() {
	app.ConnectDb("./urlshortner.db")

	//start a go routine in separate thread to discard all failed login attempt records after one hour
	app.FailedAttempts = make(map[string]*FailedAttempt)
	app.CleanFailedAttempts(time.Minute * 60)

	//initialize session with a secret to encrypt cookies
	gob.Register(User{})
	app.InitSessionStore()
}

//connects to sqlite3 database file urlshortner.db, if it's not already present this function creates it
//This function fails with os.Exit(1) if it couldn't connect to db
func (app *App) ConnectDb(name string) {
	//var db *sqlx.DB

	// from a pre-existing sql.DB; note the required driverName
	c, err := sql.Open("sqlite3", name)

	if err != nil {
		log.Fatalln("%s", err)
	}

	app.Db = sqlx.NewDb(c, "sqlite3")

	// force a connection and test that it worked
	err = app.Db.Ping()
	if err != nil {
		log.Fatalln("Could not connect database. %s", err)
	}
	err = app.createDb()
	if err != nil {
		log.Fatalf("\nError %s\n", err)
	}
}

//compares two slices of strings and returns number of values in expected[] slice which are not in given[] slice
func diffSlices(expected []string, given []string) int {
	c := 0
	for _, t := range given {
		for _, e := range expected {
			if e == t {
				c += 1
			}
		}
	}
	return len(expected) - c
}

//Creates database schema for our web app
func (app *App) createDb() error {
	schema := `
	BEGIN;
	CREATE TABLE IF NOT EXISTS users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    email TEXT,
	    password TEXT,
	    ondate INTEGER
    );
	CREATE INDEX IF NOT EXISTS id_email ON users (email);

	CREATE TABLE IF NOT EXISTS  urls (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    url TEXT,
	    userid INTEGER NULL,
	    ondate INTEGER
    );
	CREATE INDEX IF NOT EXISTS id_userid ON urls (userid);

	CREATE TABLE IF NOT EXISTS hits (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    referrer TEXT,
	    ip TEXT,
	    urlid INTEGER,
	    ondate INTEGER,
	    FOREIGN KEY(urlId) REFERENCES urls(id)
    );
	CREATE INDEX IF NOT EXISTS id_urlid ON hits (urlid);
	INSERT INTO URLS (id) VALUES(1000);
	DELETE FROM URLS WHERE id = 1000;
	COMMIT;
	`
	// execute a query on the server
	_, err := app.Db.Exec(schema)

	return err
}

//initializes sessions for app
func (app *App) InitSessionStore() {
	app.SessionStore = sessions.NewCookieStore([]byte("hello here's my super secret secret"))
}

func (app *App) GetSession(r *http.Request) *sessions.Session {
	session, _ := app.SessionStore.Get(r, "urlshortner-session")
	return session
}

//records a failed attempt to login via an email
func (app *App) RecordFailedAttempt(email string) {
	if _, ok := app.FailedAttempts[email]; !ok {
		app.FailedAttempts[email] = &FailedAttempt{
			Count: 1,
			Last:  time.Now().Unix(),
		}
	} else {
		app.FailedAttempts[email].Last = time.Now().Unix()
		app.FailedAttempts[email].Count += 1
	}
}

//checks if some email has previous failed attempts within 30 seconds
func (app *App) HasFailedAttempts(email string) bool {
	if a, ok := app.FailedAttempts[email]; ok {
		now := time.Now().Unix()
		if now-a.Last < 30 && a.Count >= 5 {
			return true
		} else if now-a.Last > 30 {
			delete(app.FailedAttempts, email)
		}
	}
	return false
}

//Starts a ticker and go routine and cleans up expired failed attempt counts after given duration
func (app *App) CleanFailedAttempts(duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now().Unix()
				for k, v := range app.FailedAttempts {
					if now-v.Last > 30 {
						delete(app.FailedAttempts, k)
					}
				}
			}
		}
	}()
}

func GetJson(v interface{}) []byte {
	j, err := json.Marshal(v)
	if err != nil {
		log.Printf("Couldn't encode json for %v. Error: %s", j, err)
	}
	return j
}
