package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"net/http"
)

const (
	salt = ")()(*()&*IO"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid Parameter"))
		return
	}
	enc_password := util.Sha1([]byte(password + salt))
	suc := dblayer.UserSignUp(username, enc_password)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))

	}
}
