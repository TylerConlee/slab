package datastore

import "time"

var id = 0

// Check tables

// If tables don't exist, create tables

// Save activity
func SaveActivity(user string, activityType string) error {
	if activityType == "set" {
		if id != 0 {
			_, err := db.Query("UPDATE activities SET ended_at = $1 WHERE id = $2", time.Now(), id)
			id = 0
			return err
		}
		err := db.QueryRow("INSERT INTO activities(slack_id, type, started_at) VALUES ($1,$2,$3) RETURNING id", user, activityType, time.Now()).Scan(&id)
		return err
	} else if activityType == "unset" {
		if id != 0 {
			_, err := db.Query("UPDATE activities SET ended_at = $1 WHERE id = $2", time.Now(), id)
			id = 0
			return err
		}
		return nil
	}
	_, err := db.Query("INSERT INTO activities(slack_id, type, started_at, ended_at) VALUES ($1,$2,$3,$4)", user, activityType, time.Now(), time.Now())
	return err

}

// Check if open activity

// Save User

func CreateActivitiesTable() {
	const activities = `
	CREATE TABLE IF NOT EXISTS activities (
		id serial PRIMARY KEY,
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

func CreateUsersTable() {
	const users = `
	CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		slack_id text NOT NULL,
		name text NOT NULL
	)`

	// Exec executes a query without returning any rows.
	if _, err := db.Exec(users); err != nil {
		log.Error(
			"Users table creation query failed", map[string]interface{}{
				"module": "datastore",
				"query":  users,
			})
		return
	}

	return
}
