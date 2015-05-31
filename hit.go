package urlshortner

import (
	"time"
)

// A hit is recorded when someone browses a shortened url
// Hit struct saves records for hits table
type Hit struct {
	Id       int64
	Referrer string
	Ip       string
	Urlid    int64
	Ondate   int64
}

//Saves a hit
func (h *Hit) Save() error {
	h.Ondate = time.Now().Unix()
	query := "INSERT INTO hits (ip, referrer, urlid, ondate) values (:ip , :referrer, :urlid, :ondate)"
	r, err := app.Db.NamedExec(query, h)
	if err != nil {
		return err
	}
	h.Id, _ = r.LastInsertId()
	return nil
}
