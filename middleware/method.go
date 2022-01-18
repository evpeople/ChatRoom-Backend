package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func Method(m string) Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			if r.Method != m {
				logrus.Trace(r.URL)
				logrus.Trace("method is" + r.Method)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			logrus.Trace("method is" + r.Method)

			f(w, r)
		}
	}
}
