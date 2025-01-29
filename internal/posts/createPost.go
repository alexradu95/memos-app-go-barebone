package posts

import (
	"database/sql"
	"time"
)

type Post struct {
	Id        int64  `db:"id"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	AccountId int64  `db:"account_id"`
}

func (p Post) FormattedDate() string {
	t, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return p.CreatedAt // Return original if parsing fails
	}
	return t.Format("January 2, 2006")
}

func CreatePost(db *sql.DB, newPost Post) (Post, error) {
	query := `INSERT INTO posts 
						(content, created_at, updated_at, account_id) 
						VALUES (?, ?, ?, ?)
						RETURNING id`

	err := db.QueryRow(
		query,
		newPost.Content,
		newPost.CreatedAt,
		newPost.UpdatedAt,
		newPost.AccountId,
	).Scan(&newPost.Id)
	if err != nil {
		return newPost, err
	}

	return newPost, nil
}
