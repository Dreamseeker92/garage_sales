package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
)

func Open() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")
	
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("sales", "sales"),
		Host:     "localhost",
		Path:     "sales_db",
		RawQuery: q.Encode(),
	}
	
	return sqlx.Open("postgres", u.String())
}