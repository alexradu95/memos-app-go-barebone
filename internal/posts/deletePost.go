package posts

import (
	"database/sql"
)

func DeletePost(db *sql.DB, postId int64) error {
	query := `DELETE FROM posts WHERE id = ?`

	_, err := db.Query(query, postId)
	if err != nil {
		return err
	}

	return nil
}
