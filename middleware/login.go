package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       uint64 `json:"id"`
}

var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

var user2 = User{
	ID:       2,
	Username: "evpeople",
	Password: "password",
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "login.html")
	} else {
		logrus.Trace("in login")
		var u User
		//http转json
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			logrus.Debug("json wrong")
			logrus.Debug(err)
		}

		if user.Username != u.Username || user.Password != u.Password {
			logrus.Debug("wrong user" + u.Username + "/n" + u.Password)
			// return
		}
		token, err := createToken(u.ID, u.Username)
		if err != nil {
			return
		}

		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}
}
