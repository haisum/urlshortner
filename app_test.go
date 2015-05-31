package urlshortner

import (
	"os"
	"testing"
)

var app = new(App)

//remove test db file before and after running tests .
func TestMain(m *testing.M) {
	os.Remove("urlshortner_test.db")
	app.ConnectDb("urlshortner_test.db")

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
