package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v4"
)

// SECRET_HASH を計算する関数
func calculateSecretHash(clientSecret, username, clientID string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Cognitoに新規ユーザーを登録する関数
func Signup(
	clientID string,
	clientSecret string,
	email string,
	password string,
) (*cognitoidentityprovider.SignUpOutput, error) {

	// AWSセッションを作成
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	svc := cognitoidentityprovider.New(sess)

	// クライアントシークレットを用いてSecretHashを計算
	secretHash := calculateSecretHash(clientSecret, email, clientID)

	// サインアップ用のリクエストを作成
	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientID),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	// Cognitoにサインアップリクエストを送信
	result, err := svc.SignUp(input)
	if err != nil {
		return nil, err
	}

	// 結果を返却
	return result, nil
}

// メールに送信された確認コードを使って、Cognitoでユーザーを有効化する関数
func ConfirmCode(
	clientID string,
	clientSecret string,
	email string,
	confirmationCode string,
) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	// AWSセッションを初期化（リージョンは東京）
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))

	// Cognitoクライアントを作成
	svc := cognitoidentityprovider.New(sess)

	// シークレットハッシュを計算
	secretHash := calculateSecretHash(clientSecret, email, clientID)

	// 確認コードとユーザー情報を設定
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(secretHash),
	}

	// サインアップ確認を実行
	result, err := svc.ConfirmSignUp(input)
	if err != nil {
		return nil, err
	}

	// 結果を返す
	return result, nil
}

func Login(clientID string, clientSecret string, email string, password string) (*cognitoidentityprovider.AuthenticationResultType, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
	svc := cognitoidentityprovider.New(sess)

	secretHash := calculateSecretHash(clientSecret, email, clientID)

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(clientID),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(email),
			"PASSWORD":    aws.String(password),
			"SECRET_HASH": aws.String(secretHash),
		},
	}

	resp, err := svc.InitiateAuth(input)
	if err != nil {
		return nil, err
	}

	return resp.AuthenticationResult, nil
}

func CheckAccessToken(accessToken string) bool {
	jwksURL := "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_gH7L1MXyK/.well-known/jwks.json"

	options := keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Printf("JWKs refresh error: %v", err)
		},
	}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		log.Printf("JWK取得失敗: %v", err)
		return false
	}

	token, err := jwt.Parse(accessToken, jwks.Keyfunc)
	if err != nil {
		log.Printf("JWT検証エラー: %v", err)
		return false
	}

	if !token.Valid {
		log.Println("トークンは無効です")
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("クレームの取得に失敗しました")
		return false
	}

	log.Printf("スコープ: %v", claims["scope"])

	return true
}

// 認証チェック関数：成功すれば true、失敗すれば false を返す
func IsAuthenticated(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	const prefix = "Bearer "
	accessToken := strings.TrimPrefix(authHeader, prefix)

	return CheckAccessToken(accessToken)
}
