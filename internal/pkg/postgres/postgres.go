package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Postgreser interface {
}

type Postgres struct {
	ConnectString string
	SqlClient     *sql.DB
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

// NewPostgresConnection - returns a Postgres connection and its conenct string for reference.  This can then be passed down
// through the service so that multiple connections do need to be made unncessarily.
func NewPostgresConnection(c PostgresConfig, l *zerolog.Logger) (db *Postgres, err error) {
	cString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.SSLMode)
	m := &Postgres{
		ConnectString: cString,
	}
	m.SqlClient, err = sql.Open("postgres", cString)
	if err != nil {
		return nil, err
	}
	err = m.SqlClient.Ping()
	if err != nil {
		l.Error().Msgf("unable to connect to postgres")
		return
	}
	l.Error().Msgf("connection to Postgres works!")

	return m, nil
}
