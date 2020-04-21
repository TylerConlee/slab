package datastore

import (
	"database/sql"
	"fmt"

	// pq is recommended to be loaded in blank
	_ "github.com/lib/pq"
	c "github.com/tylerconlee/slab/config"
	l "github.com/tylerconlee/slab/log"
)

var (
	log = l.Log
	db  *sql.DB
)

// PGConnect uses the configuration passed from the config file to connect to
// Postgres and ensure that the table is created properly.
func Connect(cfg c.Config) {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Error("Error encountered attempting to connect to Postgres.", map[string]interface{}{
			"error": err,
		})
	}
	err = db.Ping()
	if err != nil {
		log.Error("Error encountered attempting to connect to Postgres.", map[string]interface{}{
			"error": err,
		})
	}
	log.Info("Postgres connected.", map[string]interface{}{
		"module": "datastore",
	})
	CreateActivitiesTable()
	CreateTagsTable()
	CreateTriagerTable()
}
