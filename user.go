package urlshortner

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
)

// User type holds data for a single user record in database
// Id is auto incremented so isn't needed when calling Save()
// Password must be at least 6 chars long and is automatically hashed via bcrypt before saving in database
// Email should be valid email address
// Ondate is UNIX timestamp at time of saving
type User struct {
	Id       int64
	Email    string
	Password string
	Ondate   int64
}

// Save method should be called on User object with valid Email address and password > 6 chars filled in.
// It validates input. Checks if email already exists.
// Sets Ondate to current unix timestamp.
// Calculates hash for password and saves user in users table
//
// Takes db.Con in argument
func (u *User) Save(con *sqlx.DB) error {
	password := []byte(u.Password)
	// Hashing the password with the cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		return err
	}
	//we need this for tests
	tempPass := u.Password
	u.Password = string(hashedPassword)

	u.Ondate = time.Now().Unix()

	query := "INSERT INTO users (email, password, onDate) values (:email , :password, :ondate)"
	r, err := con.NamedExec(query, u)
	if err != nil {
		return err
	}
	u.Id, _ = r.LastInsertId()
	u.Password = tempPass

	return nil
}

// Validates a user record for password length and email
// returns nil if valid record and
// returns []error if there were errors encountered
func (u User) Validate(con *sqlx.DB) []error {
	errs := make([]error, 0)
	emailPattern := regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
	user := User{
		Email: u.Email,
	}
	if len(strings.TrimSpace(u.Email)) == 0 {
		errs = append(errs, errors.New("Email must not be empty"))
	} else if !emailPattern.MatchString(u.Email) {
		errs = append(errs, errors.New("Not a valid Email address."))
	} else if user.EmailGet(con); user.Id != 0 {
		errs = append(errs, errors.New("Email already registered."))
	}

	if len(u.Email) > 50 {
		errs = append(errs, errors.New("Email can't be longer than 50 chars."))
	}

	if len(u.Password) < 6 {
		errs = append(errs, errors.New("Password must be at least 6 chars long."))
	}
	if len(u.Password) > 30 {
		errs = append(errs, errors.New("Password can't be more than 30 chars long."))
	}
	if len(errs) > 0 {
		return errs
	} else {
		return nil
	}
}

// Email Get method should be called on User object with Email field set to something.
// It checks table users for record with email = u.Email, fills u with record on success and returns nil
// or returns error on failure
//
// Takes db.Con in argument
func (u *User) EmailGet(con *sqlx.DB) error {
	st, err := con.PrepareNamed("SELECT * FROM users where email = :email")
	if err != nil {
		return err
	}
	err = st.Get(u, u)
	return err
}
