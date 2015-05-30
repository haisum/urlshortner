package urlshortner

import (
	"testing"
	"time"
)

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
