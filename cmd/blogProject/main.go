package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"blogProject/internal"
)

func main() {
	internal.InitializeDB("root:mdp2nlucK.#@tcp(127.0.0.1:3306)/blogproject")

	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/user", CreateUserHandler)
	http.HandleFunc("/user/profile", UserProfileHandler)
	http.HandleFunc("/articles", ArticlesHandler)
	http.HandleFunc("/articles/create", CreateArticleHandler)
	http.HandleFunc("/articles/{articleID}/comment", AddCommentHandler)
	http.HandleFunc("/articles/{articleID}/like", LikeArticleHandler)
	http.HandleFunc("/articles/{articleID}/dislike", DislikeArticleHandler)
	http.HandleFunc("/articles/{articleID}", DeleteArticleHandler)

	port := 8080
	fmt.Printf("Serveur démarré sur le port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

// HomeHandler gère l'endpoint /home
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	articles, err := internal.HomeArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// LoginHandler gère l'endpoint /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	email, exists := request["email"]
	if !exists || email == "" {
		http.Error(w, "Champ email manquant", http.StatusBadRequest)
		return
	}

	err = internal.LoginUser(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateUserHandler gère l'endpoint /user (POST)
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var user internal.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	err = internal.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UserProfileHandler gère l'endpoint /user/profile (GET)
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	user, articles, err := internal.GetUserProfile(email[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		User     internal.User      `json:"user"`
		Articles []internal.Article `json:"articles"`
	}{
		User:     *user,
		Articles: articles,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// ArticlesHandler gère l'endpoint /articles (GET)
func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	articles, err := internal.GetArticlesByEmail(email[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// CreateArticleHandler gère l'endpoint /articles (POST)
func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	var article internal.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	article.UserID, err = internal.GetUserIDByEmail(email[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = internal.CreateArticle(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AddCommentHandler gère l'endpoint /articles/{articleID}/comment (POST)
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	articleID, err := strconv.Atoi(r.URL.Query().Get("articleID"))
	if err != nil {
		http.Error(w, "Paramètre articleID invalide", http.StatusBadRequest)
		return
	}

	var comment internal.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	err = internal.AddComment(articleID, email[0], &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// LikeArticleHandler gère l'endpoint /articles/{articleID}/like (POST)
func LikeArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	articleID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articles/"))
	if err != nil {
		http.Error(w, "Paramètre articleID invalide", http.StatusBadRequest)
		return
	}

	err = internal.LikeArticle(articleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DislikeArticleHandler gère l'endpoint /articles/{articleID}/dislike (POST)
func DislikeArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	articleID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articles/"))
	if err != nil {
		http.Error(w, "Paramètre articleID invalide", http.StatusBadRequest)
		return
	}

	err = internal.DislikeArticle(articleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteArticleHandler gère l'endpoint /articles/{articleID} (DELETE)
func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	email, exists := r.Header["Email"]
	if !exists || email[0] == "" {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}

	articleID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/articles/"))
	if err != nil {
		http.Error(w, "Paramètre articleID invalide", http.StatusBadRequest)
		return
	}

	err = internal.DeleteArticle(articleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
