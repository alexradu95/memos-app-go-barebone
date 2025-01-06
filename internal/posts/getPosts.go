package posts

import (
	"database/sql"
	"time"
)

type GetPostsResponse struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	Id          string `json:"id"`
	Content     string `json:"content"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
	UserId      string `json:"userId"`
}

type QueryParams struct {
	UserId     string `query:"userId"`
	SearchText string `query:"searchText"`
	DateFrom   string `query:"dateFrom"`
	DateTo     string `query:"dateTo"`
	PageNumber int    `query:"pageNumber"`
	PageSize   int    `query:"pageSize"`
}

func GetPosts(db *sql.DB, params QueryParams) []Post {
	args := []interface{}{params.UserId}
	query := `SELECT id, content, date_created AS dateCreated, date_updated AS dateUpdated, user_id AS userId
		          FROM posts
		          WHERE user_id = $1`

	if params.SearchText != "" {
		query += " AND content ILIKE $2"
		args = append(args, "%"+params.SearchText+"%")
	}

	if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err != nil {
			return []Post{}
		}
		query += " AND date_created >= $3"
		args = append(args, dateFrom)
	}

	if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err != nil {
			return []Post{}
		}
		query += " AND date_created <= $4"
		args = append(args, dateTo)
	}

	query += " ORDER BY date_created DESC"

	if params.PageNumber > 0 && params.PageSize > 0 {
		offset := (params.PageNumber - 1) * params.PageSize
		query += " LIMIT $5 OFFSET $6"
		args = append(args, params.PageSize, offset)
	}

	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return []Post{}
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.Id,
			&post.Content,
			&post.DateCreated,
			&post.DateUpdated,
			&post.UserId,
		)
		if err != nil {
			return []Post{}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return []Post{}
	}

	return posts
}

func GetPost(db *sql.DB, userId string, postId string) Post {
	query := `SELECT id, content, date_created AS dateCreated, date_updated AS dateUpdated, user_id AS userId
		          FROM posts
		          WHERE user_id = $1
		          AND id = $2`

	args := []interface{}{userId, postId}

	rows, err := db.Query(query, args...)
	if err != nil {
		return Post{}
	}

	defer rows.Close()

	var post Post

	for rows.Next() {
		err := rows.Scan(
			&post.Id,
			&post.Content,
			&post.DateCreated,
			&post.DateUpdated,
			&post.UserId,
		)
		if err != nil {
			return Post{}
		}
	}

	if err := rows.Err(); err != nil {
		return Post{}
	}

	return post
}
