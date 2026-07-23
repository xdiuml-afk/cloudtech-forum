package main

import (
	"log"
	"net/http"
	"os"

	"cloudtech-forum/handler"
	"cloudtech-forum/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// .envファイルから環境変数を読み込む
	godotenv.Load()
}

// CORSミドルウェア
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORSヘッダーを設定（必要に応じて制限することも可能）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// OPTIONSメソッドへの即時レスポンス（Preflightリクエスト対応）
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 通常処理
		next.ServeHTTP(w, r)
	})
}
func main() {

	// 環境変数からデータを取得
	apiport := os.Getenv("API_PORT")     // APIサーバのポート
	username := os.Getenv("DB_USERNAME") // DBのユーザ名
	password := os.Getenv("DB_PASSWORD") // DBのパスワード
	host := os.Getenv("DB_HOST")         // DBのホスト
	port := os.Getenv("DB_PORT")         // DBのポート
	dbname := os.Getenv("DB_NAME")       // DB名

	// データベースの接続を初期化
	err := repository.InitDB(username, password, host, port, dbname)
	if err != nil {
		log.Fatalf("データベースに接続できません: %v", err)
	}
	defer repository.CloseDB() // プログラム終了時にデータベース接続を閉じる

	// ルーティングの設定
	r := mux.NewRouter()
	r.HandleFunc("/posts", handler.CreateHandler).Methods("POST")
	r.HandleFunc("/posts", handler.IndexHandler).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}", handler.ShowHandler).Methods("GET")
	r.HandleFunc("/posts/{id:[0-9]+}", handler.UpdateHandler).Methods("PUT")
	r.HandleFunc("/posts/{id:[0-9]+}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/signup", handler.SignupHandler).Methods("POST")
	r.HandleFunc("/confirmcode", handler.ConfirmSignupHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	// CORSミドルウェアを適用
	corsRouter := enableCORS(r)

	// APIサーバを起動
	log.Println("APIサーバを起動しました。ポート: " + apiport)
	if err := http.ListenAndServe(":"+apiport, corsRouter); err != nil {
		log.Fatal(err)
	}
}
