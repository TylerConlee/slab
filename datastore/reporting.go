package datastore

import (
	"database/sql/driver"
	"time"
)

var id = 0

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// ActivityOptions allows us to pass multiple optional parameters to the LoadActivity function, including a filter for activityType and a cap on how many are loaded
type ActivityOptions struct {
	ActivityType string
	Quantity     int
}

// Activity represents a single activity listed in the database
type Activity struct {
	SlackID      string
	SlackName    string
	ActivityType string
	StartedAt    NullTime
	EndedAt      NullTime
}

// SaveActivity takes the user data and activity type and saves it to the
// Postgres database
func SaveActivity(user string, name string, activityType string) error {
	if activityType == "set" {
		if id != 0 {
			_, err := db.Query("UPDATE activities SET ended_at = $1 WHERE id = $2", time.Now(), id)
			err = db.QueryRow("INSERT INTO activities(slack_id, slack_name, type, started_at) VALUES ($1,$2,$3, $4) RETURNING id", user, name, activityType, time.Now()).Scan(&id)
			return err
		}
		err := db.QueryRow("INSERT INTO activities(slack_id, slack_name, type, started_at) VALUES ($1,$2,$3, $4) RETURNING id", user, name, activityType, time.Now()).Scan(&id)
		return err
	} else if activityType == "unset" {
		if id != 0 {
			_, err := db.Query("UPDATE activities SET ended_at = $1 WHERE id = $2", time.Now(), id)
			id = 0
			return err
		}
		return nil
	}
	_, err := db.Query("INSERT INTO activities(slack_id, slack_name, type, started_at, ended_at) VALUES ($1,$2,$3,$4, $5)", user, name, activityType, time.Now(), time.Now())
	return err

}

// LoadActivity requests a list of activity/activities from the database to use in the History command, or in reporting
func LoadActivity(opts ActivityOptions) (activities []Activity, err error) {
	if opts.Quantity == 0 {
		opts.Quantity = 30
	}
	rows, err := db.Query("select slack_id, slack_name, type, started_at, ended_at from activities ORDER BY started_at DESC LIMIT $1;", opts.Quantity)
	if err != nil {
		log.Error("Error encountered attempting to load from Postgres.", map[string]interface{}{
			"error": err,
		})
	}
	defer rows.Close()
	for rows.Next() {
		act := Activity{}
		err := rows.Scan(&act.SlackID, &act.SlackName, &act.ActivityType, &act.StartedAt, &act.EndedAt)
		if err != nil {
			log.Error("Error scanning loaded activity.", map[string]interface{}{
				"error": err,
			})
		}
		activities = append(activities, act)
	}
	return activities, err
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
