package datastore

import "time"

// CreateTriagerTable checks to see if the proper table exists, and if it
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

// CreateChannelsTable checks to see if the proper table exists, and if it
// doesn't, create one.
func CreateChannelsTable() {
	const channels = `
	CREATE TABLE IF NOT EXISTS channels (
		id serial PRIMARY KEY,
		channel_id text NOT NULL,
		updated_at timestamp,
		private boolean

	)`

	// Exec executes a query without returning any rows.
	if _, err := db.Exec(channels); err != nil {
		log.Error("Channels table creation query failed", map[string]interface{}{
			"module": "datastore",
			"error":  err,
			"query":  channels,
		})
		return
	}

	return
}

// SaveTriager saves a new triager into the database
func SaveTriager(data string) error {
	log.Info("Preparing triager for database", map[string]interface{}{
		"module": "datastore",
		"data":   data,
	})
	err := db.QueryRow("INSERT INTO triager(userid,  updated_at) VALUES ($1, $2) RETURNING id", data, time.Now()).Scan(&id)
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

// SaveChannels saves channels into the database and updates the existing
// record if one exists
func SaveChannels(data map[string]interface{}) error {
	log.Info("Preparing channels for database", map[string]interface{}{
		"module": "datastore",
		"data":   data,
	})
	err := db.QueryRow("INSERT INTO channels(channels, private, updated_at) VALUES ($1, $2, $3) WHERE id = 1 ON CONFLICT (id) DO UPDATE SET channels = $1 RETURNING id", data["channel_id"], data["private"], time.Now()).Scan(&id)
	return err
}

// LoadChannels grabs the list of channels from the database and returns them
// in a slice of strings.
func LoadChannels(dm bool) (channels []string, err error) {
	log.Info("Requesting channels from database", map[string]interface{}{
		"module": "datastore",
	})
	rows, err := db.Query("SELECT channel_id FROM channels WHERE private = $1", dm)
	if err != nil {
		log.Error("Error grabbing database output for channels", map[string]interface{}{
			"module": "datastore",
			"error":  err,
		})
	}
	defer rows.Close()

	for rows.Next() {
		var (
			channelID string
		)

		if err := rows.Scan(&channelID); err != nil {
			log.Error("Error parsing database output for channels", map[string]interface{}{
				"module": "datastore",
				"error":  err,
			})
		}

		channels = append(channels, channelID)
	}
	return
}
