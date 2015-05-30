package urlshortner

import (
	"os"
	"testing"
	"time"
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

func TestUser(t *testing.T) {
	db := Db{
		Name: "urlshortner_test.db",
	}
	db.ConnectDb()
	con := db.Con

	data := User{
		Email:    "haisumbhatti@gmail.com",
		Password: "jasdjkjksdf",
		Ondate:   time.Now().Unix(),
	}

	for x := 0; x < 10; x++ {
		query := "INSERT INTO users (email, password, onDate) values (:email , :password, :ondate)"
		_, err := con.NamedExec(query, data)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}

	users := []User{}
	query := "SELECT * FROM users"
	err := con.Select(&users, query)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if len(users) != 10 {
		t.Fatalf("Not enough users! Total users: %d, data: %v", len(users), users)
	}
}

func TestUrl(t *testing.T) {
	db := Db{
		Name: "urlshortner_test.db",
	}
	db.ConnectDb()
	con := db.Con

	data := Url{
		Longurl: "http://github.com/mysuperlongurl",
		Userid:  1,
		Ondate:  time.Now().Unix(),
	}

	for x := 0; x < 10; x++ {
		query := "INSERT INTO urls (longurl, userid, ondate) values (:longurl , :userid, :ondate)"
		_, err := con.NamedExec(query, data)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}

	urls := []Url{}
	query := "SELECT * FROM urls"
	err := con.Select(&urls, query)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if len(urls) != 10 {
		t.Fatalf("Not enough urls! Total urls: %d, data: %v", len(urls), urls)
	}
}
