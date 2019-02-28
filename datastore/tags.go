package datastore

// CreateTagsTable checks to see if the proper table exists, and if it
// doesn't, create one.
func CreateTagsTable() {
	const tags = `
	CREATE TABLE IF NOT EXISTS tags (
		id serial PRIMARY KEY,
		tag text NOT NULL,
		user text NOT NULL,
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
