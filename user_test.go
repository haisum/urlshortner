package urlshortner

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	con := app.Db

	data := User{
		Email:    "haisumbhatti@gmail.com",
		Password: "jasdjkjksdf",
		Ondate:   time.Now().Unix(),
	}

	for x := 0; x < 10; x++ {
		query := "INSERT INTO users (email, password, onDate) values (:email , :password, :ondate)"
		_, err := con.NamedExec(query, data)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}

	users := []User{}
	query := "SELECT * FROM users"
	err := con.Select(&users, query)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if len(users) != 10 {
		t.Fatalf("Not enough users! Total users: %d, data: %v", len(users), users)
	}
}

func TestUser_Save(t *testing.T) {
	con := app.Db

	users := []User{
		{Email: "haisum1@gmail.com", Password: "helloworld"},
		{Email: "haisum2@gmail.com", Password: "helloworld2"},
		{Email: "haisum3@gmail.com", Password: "helloworld3"},
		{Email: "haisum4@gmail.com", Password: "helloworld4"},
		{Email: "haisum5@gmail.com", Password: "helloworld5"},
		{Email: "haisum6@gmail.com", Password: "hellowor6ld5"},
	}

	for k, _ := range users {
		err := users[k].Save()
		if err != nil {
			t.Fatalf("Error: %s", err)
		}
	}
	dbusers := []User{}
	query := fmt.Sprintf("SELECT * FROM users where id >= %d AND id <= %d", users[0].Id, users[len(users)-1].Id)
	err := con.Select(&dbusers, query)

	if err != nil {
		t.Fatalf("Error: %s. Query: %s", err, query)
	}

	if len(users) != len(dbusers) {
		t.Fatalf("Not enough users! Total users: %d, data: %v. Users in db: %d, data: %v", len(users), users, len(dbusers), dbusers)
	}

	for x := 0; x < len(users); x++ {
		err = bcrypt.CompareHashAndPassword([]byte(dbusers[x].Password), []byte(users[x].Password))
		if err != nil {
			t.Fatalf("Hash didn't work! Hash: %s, Plain: %s, UserId1: %d, Userid2: %d", dbusers[x].Password, users[x].Password, users[x].Id, dbusers[x].Id)
		}
	}
}

func TestUser_Validate(t *testing.T) {
	//all valid users
	users := []User{
		{Email: "haisum7@gmail.com", Password: "helloworld"},
		{Email: "haisum8@gmail.com", Password: "helloworld2"},
		{Email: "haisum9@gmail.com", Password: "helloworld3"},
		{Email: "haisum10@gmail.com", Password: "helloworld4"},
		{Email: "haisum11@gmail.com", Password: "helloworld5"},
		{Email: "haisum12@gmail.com", Password: "hellowor6ld5"},
	}

	for k, _ := range users {
		err := users[k].Validate()
		if err != nil {
			t.Fatalf("\nError: %s\n", err)
		}
	}
	var user User
	var errs []error

	user = User{
		Email:    "bademail@.c",
		Password: "ioweiorioweio",
	}
	errs = user.Validate()
	if errs == nil || errs[0].Error() != "Not a valid Email address." || len(errs) > 1 {
		t.Fatalf("\nFailed validation for invalid email\n")
	}

	user = User{
		Email:    "",
		Password: "ioweirioiwoeroiweroewr",
	}
	errs = user.Validate()
	if errs == nil || errs[0].Error() != "Email must not be empty" || len(errs) > 1 {
		t.Fatalf("\nFailed validation for empty email\n")
	}

	user = User{
		Email:    "jasjdkj@sdfsfd.com",
		Password: "ioewr",
	}
	errs = user.Validate()
	if errs == nil || errs[0].Error() != "Password must be at least 6 chars long." || len(errs) > 1 {
		t.Fatalf("\nFailed validation for short password\n")
	}

	user = User{
		Email:    "weuiuiew@sdfklsalkfklsdfklskldfklklasdfklskladfklsdfkldsf.com",
		Password: "l;wld;flsa;dl;flsdaklfklsdaklfklsdaklfklsdalfkkldsflkklsdflklkdsf",
	}
	errs = user.Validate()
	if errs == nil || errs[0].Error() != "Email can't be longer than 50 chars." || errs[1].Error() != "Password can't be more than 30 chars long." || len(errs) != 2 {
		t.Fatalf("\nFailed validation for long password and email\n")
	}

	user = User{
		Email:    "haisum@abc.com",
		Password: "helloworld1",
	}
	user.Save()

	errs = user.Validate()
	if errs == nil || errs[0].Error() != "Email already registered." || len(errs) > 1 {
		t.Fatalf("\nEmail validation for already existing mail failed. %s\n")
	}

}

func TestUser_GetEmail(t *testing.T) {
	user1 := User{
		Email:    "haisum1@abc.com",
		Password: "helloworld1",
	}
	err := user1.Save()
	if err != nil {
		t.Fatalf("Couldn't insert user %v", user1)
	}

	user2 := User{
		Email: "haisum1@abc.com",
	}

	err = user2.EmailGet()

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if user2.Id != user1.Id {
		t.Fatalf("Couldn't find user with email haisum@abc.com. User  1: %v, User 2: %v", user1, user2)
	}

}
