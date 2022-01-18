package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
			return
		}
		token, err := createToken(user.ID)
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
func createToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdxfksdmfanc") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Second * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
