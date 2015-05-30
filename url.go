package urlshortner

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

// struct Url stores user supplied url
type Url struct {
	Id     int64
	Url    string
	Userid int64
	Ondate int64
}

func (u *Url) Save(con *sqlx.DB) error {
	u.Ondate = time.Now().Unix()
	query := "INSERT INTO urls (url, userid, ondate) values (:url , :userid, :ondate)"
	r, err := con.NamedExec(query, u)
	u.Id, _ = r.LastInsertId()
	return err
}

// Validates a url record for url length
// Note: pattern isn't deliberately checked because url could be anything such as localhost or http://example.com or www.example.com or localhost.dev
// By restricting pattern we may make app less usable for some users
// returns nil if valid record and
// returns []error if there were errors encountered
func (u Url) Validate() []error {
	errs := make([]error, 0)
	if len(strings.TrimSpace(u.Url)) == 0 {
		errs = append(errs, errors.New("Url must not be empty"))
	}
	if len(u.Url) > 300 {
		errs = append(errs, errors.New("Url can't be longer than 300 chars."))
	}

	if len(errs) > 0 {
		return errs
	} else {
		return nil
	}
}

//Gets a url from id
func (u *Url) Get(con *sqlx.DB) error {
	st, err := con.PrepareNamed("SELECT * FROM urls where id = :id")
	if err != nil {
		return err
	}
	err = st.Get(u, u)
	return err
}

// Gets urls of a user from offset to limit
func GetAllUrls(con *sqlx.DB, userId int, offset int, limit int) ([]Url, error) {
	var urls = []Url{}
	st, err := con.PrepareNamed("SELECT * FROM urls where userid = :userid LIMIT :offset, :limit")
	if err != nil {
		return nil, err
	}
	err = st.Select(&urls, struct {
		Userid int
		Offset int
		Limit  int
	}{
		userId,
		offset,
		limit,
	})
	return urls, err
}
