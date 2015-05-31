package urlshortner

import (
	"encoding/gob"
	"os"
	"testing"
	"time"
)

//remove test db file before and after running tests .
func TestMain(m *testing.M) {
	os.Remove("urlshortner_test.db")

	app.ConnectDb("./urlshortner_test.db")
	//start a go routine in separate thread to discard all failed login attempt records after one hour
	app.FailedAttempts = make(map[string]*FailedAttempt)
	app.CleanFailedAttempts(time.Minute * 60)

	//initialize session with a secret to encrypt cookies
	gob.Register(User{})
	app.InitSessionStore()

	r := m.Run()
	os.Remove("urlshortner_test.db")
	os.Exit(r)
}

func TestDb_ConnectDb(t *testing.T) {
	err := app.Db.Ping()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
