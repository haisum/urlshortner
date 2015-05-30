package urlshortner

import (
	"fmt"
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
		Url:    "http://github.com/mysuperlongurl",
		Userid: 1,
		Ondate: time.Now().Unix(),
	}

	for x := 0; x < 10; x++ {
		query := "INSERT INTO urls (url, userid, ondate) values (:url , :userid, :ondate)"
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

func TestUrl_Save(t *testing.T) {
	db := Db{
		Name: "urlshortner_test.db",
	}
	db.ConnectDb()
	con := db.Con

	data := []Url{
		{
			Url:    "http://github.com/mysuperlongurl",
			Userid: 1,
			Ondate: time.Now().Unix(),
		},
		{
			Url:    "htt",
			Userid: 0,
		},
		{
			Url: "klklkllkkl",
		},
	}

	for x := 0; x < len(data); x++ {
		err := data[x].Save(con)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}

	urls := []Url{}
	query := fmt.Sprintf("SELECT * FROM urls where id >= %d and id <= %d", data[0].Id, data[len(data)-1].Id)
	err := con.Select(&urls, query)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if len(urls) != len(data) {
		t.Fatalf("Not enough urls! Total urls: %d, data: %v", len(urls), urls)
	}
}

func TestUrl_Get(t *testing.T) {
	db := Db{
		Name: "urlshortner_test.db",
	}
	db.ConnectDb()
	con := db.Con

	url1 := Url{
		Url: "localhost",
	}

	url1.Save(con)

	url2 := Url{
		Id: url1.Id,
	}

	url2.Get(con)

	if url2.Url != url1.Url {
		t.Fatalf("Failed url get via id")
	}
}

func TestGetAllUrls(t *testing.T) {
	db := Db{
		Name: "urlshortner_test.db",
	}
	db.ConnectDb()
	con := db.Con

	data := []Url{
		{
			Url:    "http://github.com/mysuperlongurl",
			Userid: 1,
			Ondate: time.Now().Unix(),
		},
		{
			Url:    "htt",
			Userid: 1,
		},
		{
			Url: "klklkllkkl",
		},
		{
			Url:    "http://github.com/mysuperlongurl",
			Userid: 2,
			Ondate: time.Now().Unix(),
		},
		{
			Url:    "htt",
			Userid: 2,
		},
		{
			Url:    "klklkllkkl",
			Userid: 2,
		},
	}

	for x := 0; x < len(data); x++ {
		err := data[x].Save(con)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}
	urls, err := GetAllUrls(con, 2, 0, 10)

	if err != nil {
		t.Fatalf("%s", err)
	}

	if urls[0].Url != "http://github.com/mysuperlongurl" || urls[1].Url != "htt" || urls[2].Url != "klklkllkkl" {
		t.Fatalf("Didn't get expected data. %v", urls)
	}
}
