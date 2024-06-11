package auth

import (
	"errors"
	"os"
	"strconv"
)

var (
	jwtSecretKey     string // シークレットキー
	jwtTokenLifeTime int    // 有効期限
)

// .envから取得されたenvJWTのさまざまな情報を読み込む
func loadEnvJwt() (string, int, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY") // 環境変数からシークレットキー(署名鍵)を取得
	if secretKey == "" {
		return "", 0, errors.New("JWT_SECRET_KEY is not set")
	}
	tokenLifeTime, err := strconv.Atoi(os.Getenv("JWT_TOKEN_LIFETIME")) // トークンの有効期限
	if err != nil {
		return "", 0, err
	}

	return secretKey, tokenLifeTime, nil
}

// envから取得したデータを変数に設定
func setJwtConf(secretKey string, tokenLifeTime int) {
	jwtSecretKey = secretKey
	jwtTokenLifeTime = tokenLifeTime
}

// 認証関連の初期化
func InitAuth() error {
	// .envから取得されたenvJWTのさまざまな情報を読み込む
	secretKey, tokenLifeTime, err := loadEnvJwt()
	if err != nil {
		return err
	}

	// 変数に設定
	setJwtConf(secretKey, tokenLifeTime)

	return nil
}
