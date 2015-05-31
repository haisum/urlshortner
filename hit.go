package urlshortner

import (
	"github.com/jmoiron/sqlx"
	"time"
)

// A hit is recorded when someone browses a shortened url
// Hit struct saves records for hits table
type Hit struct {
	Id       int64
	Referrer string
	Ip       string
	Urlid    int
	Ondate   int64
}

//Saves a hit
func (h *Hit) Save(con *sqlx.DB) error {
	h.Ondate = time.Now().Unix()
	query := "INSERT INTO hits (ip, referrer, urlid, ondate) values (:ip , :referrer, :urlid, :ondate)"
	r, err := con.NamedExec(query, h)
	if err != nil {
		return err
	}
	h.Id, _ = r.LastInsertId()
	return nil
}
