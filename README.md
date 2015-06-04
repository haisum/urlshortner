Url Shortner written in google golang
==============================

This is simple web app for url shortening. 

[![Build Status](https://travis-ci.org/haisum/urlshortner.svg?branch=master)](https://travis-ci.org/haisum/urlshortner)

[![Coverage Status](https://coveralls.io/repos/haisum/urlshortner/badge.svg?branch=master)](https://coveralls.io/r/haisum/urlshortner?branch=master)


Code Documentation
------------

Source code documentation is available at http://godoc.org/github.com/haisum/urlshortner

Installation
--------------

####Manual compilation

```
go get github.com/haisum/urlshortner
cd src/github.com/haisum/urlshortner
go build urlshortner.go
./urlshortner
```

####Download pre-compiled binaries

Compiled binaries for windows and linux are available at https://drive.google.com/file/d/0B9oMRwzY0tXsc2tJenZWb29Qbjg/view?usp=sharing

**Note:** Windows binaries are not tested and may not work. I recommend running this project on linux. If windows binaries didn't work for you, please file an issue.


Features
-----------

- Thanks to golang, no need to download anything other than program binary itself making web app very portable and free  of deployment nightmares. App comes with database drivers and http server.
- Uses sqlite3 for database which can easily be replaced by any other database by making changes on couple of lines in app.go file.
- Recaptcha for users not logged in.
- Url shortening by using base36 encoding.
- User registeration.
- Saves urls for registered users.
- Stats for urls redirected.
- Beautiful, ajax based intuitive interface.
- Smart pagination of urls those logged in user shortened.
- Friendly and easy to read time formatting
- Brute force protection on login by restricting only 5 failed attempts per email per 30 seconds.
- Sql injection prevention by binding query params.
- XSS prevention by using html encoded template variables.
- Passwords hashed with bcrypt and unique salt per user
- Proper separation of concerns: handlers.go has business logic. app.go serves as global app controller. hit.go, url.go and user.go are database models. templates folder is for views and static folder for static resources.

Missing Features
--------------

- Registered users can shorten as many links as they want, it should be restricted to a limit such as 300 per day.
- CSRF protection is not coded.
- Javascript shall be more sane instead of only jQuery + Handlebars, we should use Angular or Backbone to arrange code properly.
- Pagination has some glitches.
- Stats such as IP, Referrer etc are recorded but not yet displayed, they could make good graphs.
- More user info such as name and avatar could be recorded and displayed.
- Same urls can be submitted again and again, this can defnitely be reduced to one record per url.

Screenshots
----------------

![Imgur](http://i.imgur.com/6Naa3FH.png)

![Imgur](http://i.imgur.com/uDWqdoS.png)

![Imgur](http://i.imgur.com/KBQFRNR.png)

![Imgur](http://i.imgur.com/PGTi9uz.png)

![Imgur](http://i.imgur.com/ps7eGnr.png)

![Imgur](http://i.imgur.com/9RXCI41.png)

![Imgur](http://i.imgur.com/Xv4WZwk.png)

![Imgur](http://i.imgur.com/wNs2Ocn.png)