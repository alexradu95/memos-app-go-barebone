package posts

import (
	"database/sql"
)

func DeletePost(db *sql.DB, postId string) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := db.Exec(query, postId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}

	return nil
}
