package functions

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/gofrs/uuid"
)

var taps int

func ShortLink(database *sql.DB, w http.ResponseWriter, r *http.Request) {
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
