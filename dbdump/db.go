package dbdump

import (
	"database/sql"
	"fmt"
)

var Databases = map[string]NewDB{}

type DBConfig struct {
	Database string
	Host     string
	Port     int
	Username string
	Password string
	DSN      string
}

func RegisterDB(name string, newDB NewDB) {
	Databases[name] = newDB
}

func DBNames() (d []string) {
	for k, _ := range Databases {
		d = append(d, k)
	}
	return d
}

type NewDB func(config DBConfig) (*sql.DB, error)

func init() {
	RegisterDB("mysql", NewMySQL)
}

func NewMySQL(config DBConfig) (*sql.DB, error) {
	dsn := config.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.Username, config.Password, config.Host, config.Port, config.Database)
	}
	return sql.Open("mysql", dsn)
}
