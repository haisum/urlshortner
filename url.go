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
	re := regexp.MustCompile(`^((ftp|http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|((www\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
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

// Gets urls of a user from offset to limit
func GetAllUrls(userId int, offset int, limit int) ([]Url, error) {
	var urls = []Url{}
	st, err := app.Db.PrepareNamed("SELECT * FROM urls where userid = :userid LIMIT :offset, :limit")
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

// Converts a base36 string to int representation
func UrlStringToId(s string) (int64, error) {
	return strconv.ParseInt(s, 36, 64)
}

// Converts an int to base36 string
func IdToUrlString(id int64) string {
	return strconv.FormatInt(id, 36)
}
