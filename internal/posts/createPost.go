package posts

import (
	"database/sql"
)

func CreatePost(db *sql.DB, newPost Post) (Post, error) {
	query := `INSERT INTO posts (id, content, date_created, date_updated, user_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(
		query,
		newPost.Id,
		newPost.Content,
		newPost.DateCreated,
		newPost.DateCreated,
		newPost.UserId,
	)
	if err != nil {
		return newPost, err
	}

	return newPost, nil
}
