package urlshortner

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Url struct {
	Id      int64
	Longurl string
	Userid  int64
	Ondate  int64
}

func (u *Url) Save(con *sqlx.DB) error {
	u.Ondate = time.Now().Unix()
	query := "INSERT INTO urls (longurl, userid, ondate) values (:longurl , :userid, :ondate)"
	r, err := con.NamedExec(query, u)
	u.Id, _ = r.LastInsertId()
	return err
}

func (u *Url) IdGet(con *sqlx.DB) error {
	query := "SELECT * FROM urls where id = :id"
	_, err := con.NamedExec(query, u)
	return err
}
