package functions

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

func Registerpage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("tmp/register.html")
	tmp.Execute(w, nil)
}
func Handleregister(database *sql.DB, w http.ResponseWriter, r *http.Request) {
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
