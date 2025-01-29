package posts

import (
	"database/sql"
)

type Post struct {
	Id        int    `db:"id"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	AccountId int    `db:"account_id"`
}

func CreatePost(db *sql.DB, newPost Post) (Post, error) {
	query := `INSERT INTO posts 
						(id, content, created_at, updated_at, account_id) 
						VALUES ($1, $2, $3, $4, $5)
					`
	_, err := db.Exec(
		query,
		newPost.Id,
		newPost.Content,
		newPost.CreatedAt,
		newPost.UpdatedAt,
		newPost.AccountId,
	)
	if err != nil {
		return newPost, err
	}

	return newPost, nil
}
