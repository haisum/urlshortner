package urlshortner

import (
	"errors"
	"regexp"
	"strconv"
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

func (u *Url) Save() error {
	u.Ondate = time.Now().Unix()
	re := regexp.MustCompile(`^(ftp|http|https):\/\/`)
	if !re.MatchString(u.Url) {
		u.Url = "http://" + u.Url
	}
	query := "INSERT INTO urls (url, userid, ondate) values (:url , :userid, :ondate)"
	r, err := app.Db.NamedExec(query, u)
	if err != nil {
		return err
	}
	u.Id, err = r.LastInsertId()
	return err
}

// Validates a url record for url length and pattern
// returns nil if valid record and
// returns []error if there were errors encountered
func (u Url) Validate() []error {
	errs := make([]error, 0)
	re := regexp.MustCompile(`^((ftp|http|https):\/\/)?([a-zA-Z0-9]+(\.[a-zA-Z0-9]+)+.*)$`)
	if len(strings.TrimSpace(u.Url)) == 0 {
		errs = append(errs, errors.New("Url must not be empty"))
	} else if !re.MatchString(u.Url) {
		errs = append(errs, errors.New(u.Url+" is not a valid url."))
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
func (u *Url) Get() error {
	st, err := app.Db.PrepareNamed("SELECT * FROM urls where id = :id")
	if err != nil {
		return err
	}
	err = st.Get(u, u)
	return err
}

//Get all hits for a url
func (u *Url) GetHits() ([]Hit, error) {
	st, err := app.Db.PrepareNamed("SELECT referrer, ip, ondate FROM hits where urlid = :id")
	if err != nil {
		return nil, err
	}
	hits := []Hit{}
	err = st.Select(&hits, u)
	return hits, err
}

// Gets urls of a user from offset to limit
func GetAllUrls(userId int64, offset int, limit int) ([]Url, error) {
	var urls = []Url{}
	st, err := app.Db.PrepareNamed("SELECT * FROM urls where userid = :userid LIMIT :offset, :limit")
	if err != nil {
		return nil, err
	}
	err = st.Select(&urls, struct {
		Userid int64
		Offset int
		Limit  int
	}{
		userId,
		offset,
		limit,
	})
	return urls, err
}

// Gets urls of a user from offset to limit
func GetTotalUrls(userId int64) (int64, error) {
	count := int64(0)
	st, err := app.Db.PrepareNamed("SELECT count(*) FROM urls where userid = :userid")
	if err != nil {
		return count, err
	}
	err = st.Get(&count, struct {
		Userid int64
	}{
		userId,
	})
	return count, err
}

// Converts a base36 string to int representation
func UrlStringToId(s string) (int64, error) {
	return strconv.ParseInt(s, 36, 64)
}

// Converts an int to base36 string
func IdToUrlString(id int64) string {
	return strconv.FormatInt(id, 36)
}
