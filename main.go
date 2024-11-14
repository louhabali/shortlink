package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB
var taps int

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
	http.HandleFunc("/", index)
	http.HandleFunc("/newlink", shortLink)
	http.HandleFunc("/login", loginpage)
	http.HandleFunc("/register", registerpage)
	http.HandleFunc("/handlelogin", handlelogin)
	http.HandleFunc("/handleregister", handleregister)
	http.HandleFunc("/logout", logout)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8088", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		cook, errcook := r.Cookie("user")
		tmp, err := template.ParseFiles("tmp/index.html")
		if err != nil {
			log.Fatal(err)
		}
		if errcook != nil {
			tmp.Execute(w, nil)
			return
		}
		tmp.Execute(w, cook.Value)
	} else {
		q := database.QueryRow(`SELECT LINK FROM LINK WHERE UID = ?`, r.URL.Path[1:])
		s := ""
		err := q.Scan(&s)
		if err != nil {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, s, 302)
		}
	}
}

func shortLink(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("tmp/link.html")
	uid, _ := uuid.NewV4()
	str := uid.String()[:8]
	link := r.FormValue("link")
	if link == "" {
		http.Redirect(w, r, "/", 302)
		return
	}
	taps++
	_, err := database.Exec(`
	INSERT INTO LINK (LINK,UID,TRIES) VALUES (?,?,?)`, link, str, taps)
	if err != nil {
		log.Fatal(err)
	}
	res := "http://localhost:8088/" + str
	tmp.Execute(w, res)
}
func loginpage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("tmp/login.html")
	tmp.Execute(w, nil)
}
func registerpage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("tmp/register.html")
	tmp.Execute(w, nil)
}
func handleregister(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")
	if user == "" || pass == "" {
		http.Redirect(w, r, "/register", 303)
		return
	}
	_, err := database.Exec(`
	INSERT INTO USER (USERNAME,PASSWORD) VALUES (?,?)`, user, pass)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/login", 303)
}
func handlelogin(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")
	if user == "" || pass == "" {
		http.Redirect(w, r, "/login", 303)
		return
	}
	row := database.QueryRow(`SELECT PASSWORD FROM USER WHERE USERNAME=?`, user)
	var p string
	row.Scan(&p)
	if p != pass {
		http.Redirect(w, r, "/login", 303)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user",
		Value:    user,
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.Redirect(w, r, "/", 303)
}
func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w,&http.Cookie{
		Name: "user",
		MaxAge: -1,
	})
	http.Redirect(w,r,"/login",303)
}