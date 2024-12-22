package internal

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitializeDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connexion à la base de données MySQL réussie!")
}

func HomeArticles() ([]Article, error) {
	rows, err := DB.Query("SELECT id, user_id, content, likes FROM articles ORDER BY id DESC LIMIT 5")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.UserID, &article.Content, &article.Likes)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func LoginUser(email string) error {
	_, err := DB.Exec("INSERT INTO users (email) VALUES (?) ON DUPLICATE KEY UPDATE email = email", email)
	return err
}

func CreateUser(user *User) error {
	_, err := DB.Exec("INSERT INTO users (email, firstname, lastname) VALUES (?, ?, ?)", user.Email, user.FirstName, user.LastName)
	return err
}

func GetUserProfile(email string) (*User, []Article, error) {
	var user User
	err := DB.QueryRow("SELECT id, email, firstname, lastname FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, nil, err
	}

	articles, err := GetArticlesByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	return &user, articles, nil
}

func GetArticlesByEmail(email string) ([]Article, error) {
	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT id, user_id, content, likes FROM articles WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.UserID, &article.Content, &article.Likes)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func CreateArticle(article *Article) error {
	_, err := DB.Exec("INSERT INTO articles (user_id, content) VALUES (?, ?)", article.UserID, article.Content)
	return err
}

func AddComment(articleID int, email string, comment *Comment) error {
	userID, err := GetUserIDByEmail(email)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO comments (article_id, user_id, content) VALUES (?, ?, ?)", articleID, userID, comment.Content)
	return err
}

func LikeArticle(articleID int) error {
	_, err := DB.Exec("UPDATE articles SET likes = likes + 1 WHERE id = ?", articleID)
	return err
}

func DislikeArticle(articleID int) error {
	_, err := DB.Exec("UPDATE articles SET likes = likes - 1 WHERE id = ?", articleID)
	return err
}

func DeleteArticle(articleID int) error {
	_, err := DB.Exec("DELETE FROM articles WHERE id = ?", articleID)
	return err
}

func GetUserIDByEmail(email string) (int, error) {
	var userID int
	err := DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	return userID, err
}
