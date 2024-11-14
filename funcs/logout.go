package functions

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", 303)
}
