package functions

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
)

func Index(database *sql.DB, w http.ResponseWriter, r *http.Request) {

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
