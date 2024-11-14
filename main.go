package main

import (
	"database/sql"

	f "functions/funcs"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func main() {
	var e error
	database, e = sql.Open("sqlite3", "shortLINK.db")
	if e != nil {
		log.Fatal(e)
	}
	_, err := database.Exec(`
	CREATE TABLE IF NOT EXISTS LINK (
	id INTEGER PRIMARY KEY AUTOINCREMENT ,
	LINK TEXT,
	UID TEXT UNIQUE,
	TRIES INTEGER DEFAULT 0);
	CREATE TABLE IF NOT EXISTS USER (
	id INTEGER PRIMARY KEY AUTOINCREMENT ,
	USERNAME TEXT UNIQUE,
	PASSWORD TEXT)`)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/login", f.Loginpage)
	http.HandleFunc("/register", f.Registerpage)
	http.HandleFunc("/logout", f.Logout)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f.Index(database, w, r)
	})
	http.HandleFunc("/newlink", func(w http.ResponseWriter, r *http.Request) {
		f.ShortLink(database, w, r)
	})
	http.HandleFunc("/handlelogin", func(w http.ResponseWriter, r *http.Request) {
		f.Handlelogin(database, w, r)
	})
	http.HandleFunc("/handleregister", func(w http.ResponseWriter, r *http.Request) {
		f.Handleregister(database, w, r)
	})
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8088", nil)
}
