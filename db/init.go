package db

import (
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite", "user.db")
	if err != nil {
		panic(err)
	}
	go Ping()
}
func Ping() error {
	for {
		err := DB.Ping()
		if err != nil {
			logrus.Panic("wrong db Ping")
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}
