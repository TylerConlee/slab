package datastore

import "time"

// CreateTagsTable checks to see if the proper table exists, and if it
// doesn't, create one.
func CreateTriagerTable() {
	const triager = `
	CREATE TABLE IF NOT EXISTS triager (
		id serial PRIMARY KEY,
		userid text NOT NULL,
		updated_at timestamp
	)`

	// Exec executes a query without returning any rows.
	if _, err := db.Exec(triager); err != nil {
		log.Error("Triager table creation query failed", map[string]interface{}{
			"module": "datastore",
			"error":  err,
			"query":  triager,
		})
		return
	}

	return
}

// SaveTriager saves a new triager into the database
func SaveTriager(data map[string]string) error {
	log.Info("Preparing triager for database", map[string]interface{}{
		"module": "datastore",
		"data":   data,
	})
	err := db.QueryRow("INSERT INTO triager(userid,  updated_at) VALUES ($1, $2, $3) RETURNING id", data["triager"], time.Now()).Scan(&id)
	return err
}

// LoadTriager requests the triager with the most recent updated_at value in
// the table
func LoadTriager() (triager string, err error) {
	log.Info("Requesting triager from database", map[string]interface{}{
		"module": "datastore",
	})
	row, err := db.Query("SELECT updated_at, triager FROM triager ORDER BY updated_at DESC LIMIT 1;")
	if err != nil {
		log.Error("Error grabbing database output for triager", map[string]interface{}{
			"module": "datastore",
			"error":  err,
		})
	}
	defer row.Close()

	if err = row.Scan(&triager); err != nil {
		log.Error("Error parsing database output for tags", map[string]interface{}{
			"module": "datastore",
			"error":  err,
		})
		triager = "None"
	}

	return
}
