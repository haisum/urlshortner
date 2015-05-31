package urlshortner

import (
	"fmt"
	"testing"
)

func TestHit_Save(t *testing.T) {
	con := app.Db

	data := []Hit{
		{
			Ip:       "192.123.334.3",
			Urlid:    1,
			Referrer: "http://google.com",
		},
		{
			Ip:       "192.123.334.3",
			Referrer: "http://google.com",
		},
		{
			Ip: "192.123.334.3",
		},
		{
			Referrer: "http://google.com",
		},
	}

	for x := 0; x < len(data); x++ {
		err := data[x].Save()
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}

	hits := []Hit{}
	query := fmt.Sprintf("SELECT * FROM hits where id >= %d and id <= %d", data[0].Id, data[len(data)-1].Id)
	err := con.Select(&hits, query)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if len(hits) != len(data) {
		t.Fatalf("Not enough hits! Total hits: %d, data: %v", len(hits), hits)
	}
}
