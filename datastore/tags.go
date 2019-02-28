package datastore

import "time"

// CreateTagsTable checks to see if the proper table exists, and if it
// doesn't, create one.
func CreateTagsTable() {
	const tags = `
	CREATE TABLE IF NOT EXISTS tags (
		id serial PRIMARY KEY,
		tag text NOT NULL,
		userid text NOT NULL,
		channel text NOT NULL,
		notify_type text NOT NULL,
		created_at timestamp,
		updated_at timestamp
	)`

	// Exec executes a query without returning any rows.
	if _, err := db.Exec(tags); err != nil {
		log.Error("Tag table creation query failed", map[string]interface{}{
			"module": "datastore",
			"error":  err,
			"query":  tags,
		})
		return
	}

	return
}

// SaveNewTag saves a new tag into the database
func SaveNewTag(data map[string]string) error {
	log.Info("Preparing tag for database", map[string]interface{}{
		"module": "datastore",
		"data":   data,
	})
	err := db.QueryRow("INSERT INTO tags(tag, userid, channel, notify_type, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", data["tag"], data["user"], data["channel"], data["notify_type"], time.Now()).Scan(&id)
	return err
}
