package datastore

import (
	"strconv"
	"time"
)

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
	err := db.QueryRow("INSERT INTO tags(tag, userid, channel, notify_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", data["tag"], data["user"], data["channel"], data["notify_type"], time.Now(), time.Now()).Scan(&id)
	return err
}

// SaveTagUpdate saves a tag update into the database
func SaveTagUpdate(data map[string]string) error {
	log.Info("Preparing tag to update in database", map[string]interface{}{
		"module": "datastore",
		"data":   data,
	})
	err := db.QueryRow("UPDATE tags SET tag = $1, userid = $2, channel = $3, notify_type = $4, updated_at = $5 WHERE id = $6", data["tag"], data["user"], data["channel"], data["notify_type"], time.Now(), data["id"]).Scan(&id)
	return err
}

// DeleteTag removes the tag located at the ID
func DeleteTag(data string) error {
	log.Info("Preparing to delete tag", map[string]interface{}{
		"module": "datastore",
		"tag":    data,
	})
	id, _ := strconv.Atoi(data)
	rows, err := db.Query("DELETE FROM tags WHERE id = $1", id)
	defer rows.Close()
	return err
}

// LoadTags retrieves all of the tags from the database, either to list
// or to loop through
func LoadTags() (tags []map[string]interface{}) {
	log.Info("Requesting tags from database", map[string]interface{}{
		"module": "datastore",
	})
	rows, err := db.Query("SELECT id, tag, userid, channel, notify_type, created_at, updated_at FROM tags WHERE id > 0")
	if err != nil {
		log.Error("Error grabbing database output for tags", map[string]interface{}{
			"module": "datastore",
			"error":  err,
		})
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         int64
			tag        string
			user       string
			channel    string
			notifyType string
			createdAt  string
			updatedAt  string
		)

		if err := rows.Scan(&id, &tag, &user, &channel, &notifyType, &createdAt, &updatedAt); err != nil {
			log.Error("Error parsing database output for tags", map[string]interface{}{
				"module": "datastore",
				"error":  err,
			})
		}

		tags = append(tags, map[string]interface{}{
			"id":          id,
			"tag":         tag,
			"user":        user,
			"channel":     channel,
			"notify_type": notifyType,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}
	return tags
}

// LoadTag takes an id and loads a single tag's information from the
// database
func LoadTag(id int) (loadedTag map[string]interface{}) {
	log.Info("Requesting tag from database", map[string]interface{}{
		"module": "datastore",
		"tag":    id,
	})
	rows, err := db.Query("SELECT id, tag, userid, channel, notify_type, created_at, updated_at FROM tags WHERE id = $1", id)
	if err != nil {
		log.Error("Error grabbing database output for tags", map[string]interface{}{
			"module": "datastore",
			"error":  err,
		})
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         int64
			tag        string
			user       string
			channel    string
			notifyType string
			createdAt  string
			updatedAt  string
		)

		if err := rows.Scan(&id, &tag, &user, &channel, &notifyType, &createdAt, &updatedAt); err != nil {
			log.Error("Error parsing database output for tags", map[string]interface{}{
				"module": "datastore",
				"error":  err,
			})
		}

		loadedTag = map[string]interface{}{
			"id":          id,
			"tag":         tag,
			"user":        user,
			"channel":     channel,
			"notify_type": notifyType,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		}
	}
	return loadedTag
}
