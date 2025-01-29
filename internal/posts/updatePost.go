package posts

import (
	"database/sql"
	"time"
)

func UpdatePostHandler(db *sql.DB, updatedPost Post) (Post, error) {
	query := `UPDATE posts SET content = $1, updated_at = $2 WHERE id = $3 AND`
	result, err := db.Exec(query, updatedPost.Content, time.Now(), updatedPost.Id)
	if err != nil {
		return updatedPost, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return updatedPost, err
	}

	return updatedPost, nil
}
