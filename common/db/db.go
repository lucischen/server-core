package db

import (
	"database/sql"
	"fmt"
	"time"

	// sql-driver
	_ "github.com/go-sql-driver/mysql"
)

var _ Idbs = new(DBs)

// Idbs Idbs
type Idbs interface {
	Use(name string) (c *sql.DB, ok bool)
}

// DBs DBs
type DBs struct {
	dbs map[string]*sql.DB
}

// New new a DBs
func New() *DBs {
	return &DBs{
		dbs: map[string]*sql.DB{},
	}
}

type Config struct {
	DB_MAXIDLECONNS int
	DB_MAXOPENCONNS int
	DB_MAXLIFETIME  string
	DB_READTIMEOUT  string
}

// Add Add
func (d *DBs) Add(name, DSN string, c *Config) error {
	if c.DB_MAXLIFETIME == "" || c.DB_READTIMEOUT == "" || c.DB_MAXIDLECONNS == 0 || c.DB_MAXOPENCONNS == 0 {
		return fmt.Errorf("Config cant empty. got: %v", c)
	}

	if _, ok := d.dbs[name]; !ok {
		s, err := sql.Open("mysql", DSN+"?parseTime=true&multiStatements=true&readTimeout="+c.DB_READTIMEOUT)
		if err != nil {
			return err
		}

		dur, _ := time.ParseDuration(c.DB_MAXLIFETIME)
		s.SetConnMaxLifetime(dur)
		s.SetMaxIdleConns(c.DB_MAXIDLECONNS)
		s.SetMaxOpenConns(c.DB_MAXOPENCONNS)

		d.dbs[name] = s
	}

	return nil
}

// Remove Remove
func (d *DBs) Remove(name string) {
	if _, ok := d.dbs[name]; ok {
		delete(d.dbs, name)
	}
}

// Use Use
func (d DBs) Use(name string) (c *sql.DB, ok bool) {
	c, ok = d.dbs[name]
	return
}
