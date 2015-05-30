package urlshortner

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//Db struct stores connection for web app's database and provides methods for managing and fetching data
//Name attribute defines name/location of sqlite3 database file
type Db struct {
	Con  *sqlx.DB
	Name string
}

//connects to sqlite3 database file urlshortner.db, if it's not already present this function creates it
//This function fails with os.Exit(1) if it couldn't connect to db
func (db *Db) ConnectDb() {
	//var db *sqlx.DB

	// from a pre-existing sql.DB; note the required driverName
	c, err := sql.Open("sqlite3", db.Name)

	if err != nil {
		log.Fatalln("%s", err)
	}

	db.Con = sqlx.NewDb(c, "sqlite3")

	// force a connection and test that it worked
	err = db.Con.Ping()
	if err != nil {
		log.Fatalln("Could not connect database. %s", err)
	}
	err = db.createDb()
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
func (db *Db) createDb() error {
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
	    longurl TEXT,
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

	COMMIT;
	`
	// execute a query on the server
	_, err := db.Con.Exec(schema)

	return err
}
