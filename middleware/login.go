package middleware

import (
	"encoding/json"
	"evpeople/ChatRoom/db"
	"net/http"

	"github.com/sirupsen/logrus"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       uint64 `json:"id"`
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
		raw, err := db.DB.Query("select name, password from usr where name=?", u.Username)
		if err != nil {
			logrus.Warn("wrong user")
			return
		}
		defer raw.Close()
		var name, password string
		for raw.Next() {
			if err := raw.Scan(&name, &password); err != nil {
				logrus.Warn(err)
			}
		}
		// logrus.Trace(name, password)
		if name == "" {
			logrus.Warn("User not exist")
			return
		}
		if name != u.Username || password != u.Password {
			logrus.Debug("wrong user" + u.Username + "/n" + u.Password)
			return
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
func Sign(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		logrus.Warn("wrong sign")
		logrus.Warn(err)
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte("can't sign because your fault"))

	}
	_, err = db.DB.Exec("insert into USR(NAME,PASSWORD)"+" VALUES(?,?)", u.Username, u.Password)
	if err != nil {
		logrus.Debug(err)
	}
	w.Write([]byte("success sign in"))
}
