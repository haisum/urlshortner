Url Shortner written in google golang
==============================

This is simple web app for url shortening. 

Features
-----------

- Integrated db drivers and http  server, so you just run the binary after compilation without need of installing anything else.
- Give a long url to get shortened link
- Re-captcha for non registered users
- Sqlite3 database support
- Ability to register and login so you could come back to your list of shortened urls
- Stats about clicks, geography and history of url etc

Installation
--------------

```
go get github.com/haisum/urlshortner
cd src/github.com/haisum/urlshortner
go build urlshortner.go
./urlshortner
```

Documentation
------------

Source code documentation is available at http://godoc.org/github.com/haisum/urlshortner