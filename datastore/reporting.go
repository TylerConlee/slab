package datastore

// Check tables

// If tables don't exist, create tables

// Save activity

// Check if open activity

// Save User

func CreateActivitiesTable() {
	const activities = `
	CREATE TABLE IF NOT EXISTS activities (
		id serial PRIMARY KEY,
		slack_id text NOT NULL,
		type text NOT NULL,
		started_at timestamp with time zone,
		ended_at timestamp with time zone
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
		name text NOT NULL,
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
