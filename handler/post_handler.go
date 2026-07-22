package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	model "cloudtech-forum/model"
	"cloudtech-forum/repository"

	"github.com/gorilla/mux"
)

// Createハンドラ関数
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストのBodyデータを格納するオブジェクトを定義
	var post model.Post

	// リクエストのBodyからJSONデータを取得し、post構造体に格納
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	// 登録処理を実行
	id, err := repository.CreatePost(post.Content, post.UserID)
	if err != nil {
		http.Error(w, "投稿データの登録に失敗しました", http.StatusInternalServerError)
		return
	}

	// レスポンスのBodyに、登録されたレコードのIDを指定
	response := map[string]interface{}{
		"message": "登録が成功しました",
		"id":      id,
	}

	// レスポンスを返却
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Indexハンドラ関数
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 検索処理の実行
	posts, err := repository.SearchPostAll()
	if err != nil {
		http.Error(w, "投稿データの検索に失敗しました", http.StatusInternalServerError)
		return
	}

	// ステータスコードに「200：OK」を設定
	w.WriteHeader(http.StatusOK)

	// postsデータのスライスをレスポンスとして設定
	json.NewEncoder(w).Encode(posts)
}

// Showハンドラ関数
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 検索処理の実行
	post, err := repository.SearchPost(id)
	if err != nil {
		http.Error(w, "投稿データの検索に失敗しました", http.StatusInternalServerError)
		return
	}

	// ステータスコードに「200：OK」を設定
	w.WriteHeader(http.StatusOK)

	// レスポンスのBodyに、検索した投稿データを設定
	json.NewEncoder(w).Encode(post)
}
