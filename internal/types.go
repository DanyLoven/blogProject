package internal

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type Article struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
	Likes   int    `json:"likes"`
}

type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	Content   string `json:"content"`
}
