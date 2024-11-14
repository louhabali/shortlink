package functions

import (
	"database/sql"
	"net/http"
	"text/template"
	"time"
)

func Loginpage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("tmp/login.html")
	tmp.Execute(w, nil)
}
func Handlelogin(database *sql.DB, w http.ResponseWriter, r *http.Request) {
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
