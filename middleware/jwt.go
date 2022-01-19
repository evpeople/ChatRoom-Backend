package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
func createToken(userId uint64, userName string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdxfksdmfanc") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["user_name"] = userName
	atClaims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
func extractToken(r *http.Request) string {
	tmp, err := r.Cookie("token")
	if err != nil {
		// logrus.Debug("no cookie")
		return ""
	}
	// logrus.Trace(tmp.Value)
	return tmp.Value
}
func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
func ExtractTokenMetadata(r *http.Request) (*User, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userName, ok := claims["user_name"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &User{
			Username: userName,
			ID:       userId,
		}, nil
	}
	return nil, err
}
