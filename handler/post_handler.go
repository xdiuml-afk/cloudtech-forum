package handler

import (
	"encoding/json"
	"log"
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

// Updateハンドラ関数
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// リクエストのBodyデータを格納するオブジェクトを定義
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "リクエストの形式が不正です", http.StatusBadRequest)
		return
	}

	// 更新処理の実行
	update_count, err := repository.UpdatePost(id, post.Content, post.UserID)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		http.Error(w, "更新処理に失敗しました", http.StatusInternalServerError)
		return
	}

	// 更新件数が0件の場合、404エラーを返す
	if update_count == 0 {
		http.Error(w, "更新対象のリソースが見つかりません", http.StatusNotFound)
		return
	}

	// レスポンスのBodyに更新件数をセット
	response := map[string]interface{}{
		"message":     "更新が成功しました",
		"updateCount": int(update_count),
	}

	// レスポンスを返却
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Deleteハンドラ関数
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 削除処理の実行
	delete_count, err := repository.DeletePost(id)
	if err != nil {
		http.Error(w, "削除処理に失敗しました", http.StatusInternalServerError)
		return
	}

	// 削除件数が0件の場合、404エラーを返す
	if delete_count == 0 {
		http.Error(w, "削除対象のリソースが見つかりません", http.StatusNotFound)
		return
	}

	// レスポンスのBodyに削除件数をセット
	response := map[string]interface{}{
		"message":      "削除が成功しました",
		"deletedCount": int(delete_count),
	}

	// ステータスコード200を返す
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
