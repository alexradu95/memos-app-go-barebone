package posts

import (
	"database/sql"
	"fmt"
	"time"
)

func UpdatePost(db *sql.DB, newContent string, postId int64) error {
	query := `UPDATE posts SET content = ?, updated_at = ? WHERE id = ?`
	result, err := db.Exec(query, newContent, time.Now(), postId)
	if err != nil {
		return err
	}

	// Optionally, you can check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("post with id %d not found", postId)
	}

	return nil
}
