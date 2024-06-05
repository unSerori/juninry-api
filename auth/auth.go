package auth

import (
	"juninry-api/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ユーザーidで認証トークンを生成
func GenerateToken(userUuid string) (string, error) {
	// JWTのさまざまな情報を設定
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")                         // 環境変数からシークレットキー(署名鍵)を取得
	tokenLifeTime, err := strconv.Atoi(os.Getenv("JWT_TOKEN_LIFETIME")) // トークンの有効期限
	if err != nil {
		return "", err
	}

	// uuidを作成し、
	newJti, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {
		return "", err
	}
	// テーブルを更新。
	if err := model.SaveJti(userUuid, newJti.String()); err != nil { // Userテーブルを更新
		return "", err
	}

	// クレーム部分
	claims := jwt.MapClaims{
		"id":  userUuid,        // user_uuid。
		"jti": newJti.String(), // new jti_uuid
		"exp": time.Now().Add(time.Second * time.Duration(tokenLifeTime)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)   // トークン(JWTを表す構造体)作成
	tokenString, err := token.SignedString([]byte(jwtSecretKey)) // []byte()でバイト型のスライスに変換し、SignedStringで署名。
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
