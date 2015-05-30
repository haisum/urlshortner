package urlshortner

import (
	"os"
	"testing"
)

//remove test db file before and after running tests .
func TestMain(m *testing.M) {
	os.Remove("urlshortner_test.db")
	r := m.Run()
	os.Remove("urlshortner_test.db")
	os.Exit(r)
}

func TestDb_ConnectDb(t *testing.T) {
	db := Db{
		Name: "urlshortner_test.db",
	}
	db.ConnectDb()
	con := db.Con

	err := con.Ping()

	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	db.Con.Close()
}
