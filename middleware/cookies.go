package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func CheckCookie() Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tmp, err := r.Cookie("token")
			if err != nil {
				logrus.Debug("no cookie")
				return
			}
			logrus.Trace(tmp.Value)
			f(w, r)
		}
	}
}
