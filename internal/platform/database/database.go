package database

import (
	"garagesale/internal/platform/conf"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
)

func Open(config *conf.DatabaseConfiguration) (*gorm.DB, error) {
	tlsMode := "require"
	if config.DisableTLS {
		tlsMode = "disable"
	}

	q := url.Values{}
	q.Set("sslmode", tlsMode)
	q.Set("timezone", config.Timezone)

	pgDsn := url.URL{
		Scheme:   config.Scheme,
		User:     url.UserPassword(config.Username, config.Password),
		Host:     config.Host,
		Path:     config.Name,
		RawQuery: q.Encode(),
	}

	return gorm.Open(postgres.Open(pgDsn.String()), &gorm.Config{})
}
