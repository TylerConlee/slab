package datastore

import (
	"time"
)

var id = 0

// SaveActivity takes the user data and activity type and saves it to the
// Postgres database
func SaveActivity(user string, name string, activityType string) error {
	if activityType == "set" {
		if id != 0 {
			rows, err := db.Query("UPDATE activities SET ended_at = $1 WHERE id = $2", time.Now(), id)
			err = db.QueryRow("INSERT INTO activities(slack_id, slack_name, type, started_at) VALUES ($1,$2,$3, $4) RETURNING id", user, name, activityType, time.Now()).Scan(&id)

			defer rows.Close()
			return err
		}
		err := db.QueryRow("INSERT INTO activities(slack_id, slack_name, type, started_at) VALUES ($1,$2,$3, $4) RETURNING id", user, name, activityType, time.Now()).Scan(&id)
		return err
	} else if activityType == "unset" {
		if id != 0 {
			rows, err := db.Query("UPDATE activities SET ended_at = $1 WHERE id = $2", time.Now(), id)
			id = 0

			defer rows.Close()
			return err
		}
		return nil
	}
	rows, err := db.Query("INSERT INTO activities(slack_id, slack_name, type, started_at, ended_at) VALUES ($1,$2,$3,$4, $5)", user, name, activityType, time.Now(), time.Now())
	defer rows.Close()
	return err

}

// CreateActivitiesTable checks to see if the proper table exists, and if it
// doesn't, create one.
func CreateActivitiesTable() {
	const activities = `
	CREATE TABLE IF NOT EXISTS activities (
		id serial PRIMARY KEY,
		slack_name text NOT NULL,
		slack_id text NOT NULL,
		type text NOT NULL,
		started_at timestamp,
		ended_at timestamp
	)`

	// Exec executes a query without returning any rows.
	if _, err := db.Exec(activities); err != nil {
		log.Error("Activities table creation query failed", map[string]interface{}{
			"module": "datastore",
			"query":  activities,
		})
		return
	}

	return
}
