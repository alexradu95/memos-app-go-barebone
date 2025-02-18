package posts

type Post struct {
	Id        int64  `db:"id"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	AccountId int64  `db:"account_id"`
}
