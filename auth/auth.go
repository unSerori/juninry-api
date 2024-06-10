package auth

import (
	"errors"
	"fmt"
	"juninry-api/logging"
	"juninry-api/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	jwtSecretKey     string // シークレットキー
	jwtTokenLifeTime int    // 有効期限
)

// 認証関連の初期化
func InitAuth() error {
	// .envから取得されたenvJWTのさまざまな情報を読み込む
	secretKey, tokenLifeTime, err := loadJwtEnvFile()
	if err != nil {
		return err
	}

	// 変数に設定
	setJwtConf(secretKey, tokenLifeTime)

	return nil
}

// .envから取得されたenvJWTのさまざまな情報を読み込む
func loadJwtEnvFile() (string, int, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY") // 環境変数からシークレットキー(署名鍵)を取得
	if jwtSecretKey == "" {
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

// ユーザーidで認証トークンを生成
func GenerateToken(userUuid string) (string, error) {
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
		"id":  userUuid,        // user_uuid  // クレーム内は1単語のみとしている
		"jti": newJti.String(), // new jti_uuid
		"exp": time.Now().Add(time.Second * time.Duration(jwtTokenLifeTime)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)   // トークン(JWTを表す構造体)作成
	tokenString, err := token.SignedString([]byte(jwtSecretKey)) // []byte()でバイト型のスライスに変換し、SignedStringで署名。
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// トークン解析検証
func ParseToken(tokenString string) (*jwt.Token, string, error) {
	// .envから定数をプロセスの環境変数にロード
	err := godotenv.Load(".env") // エラーを格納
	if err != nil {              // エラーがあったら
		logging.ErrorLog("Error loading .env file", err)
		panic("Error loading .env file.")
	}

	// 署名が正しければ、解析用の鍵を使う。(無名関数内で署名方法がHMACであるか確認し、HMACであれば秘密鍵を渡し、jwtトークンを解析する。)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // 署名を確認
			logging.ErrorLog(fmt.Sprintf("Unexpected signature method: %v.", token.Header["alg"]), nil)
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil // 署名が正しければJWT_SECRET_KEYをバイト配列にして返す
	})
	if err != nil {
		return nil, "", err
	}

	// 下のクレーム検証処理(:elseスコープ内)で持ち出したい値をあらかじめ宣言しておく。
	var id string // id

	// トークン自体が有効か秘密鍵を用いて確認。また、クレーム部分も取得。(トークンの署名が正しいか、有効期限内か、ブラックリストでないか。)
	claims, ok := token.Claims.(jwt.MapClaims) // MapClaimsにアサーション
	if !ok || !token.Valid {                   // 取得に失敗または検証が失敗
		return nil, "", errors.New("invalid authentication token")
	} else { // 有効な場合クレームの各要素を検証
		// idを検証
		idClaims, ok := claims["id"].(string) // goではJSONの数値は少数もカバーしたfloatで解釈される
		if !ok {
			return nil, "", errors.New("id could not be obtained from the token")
		}
		id = idClaims                           // 調整してスコープ買いに持ち出す。
		if err := model.CfmId(id); err != nil { // ユーザーに存在するか。int(id)
			return nil, "", err
		}
		// jtiを検証
		jti, ok := claims["jti"].(string)
		if !ok {
			return nil, "", errors.New("jti could not be obtained from the token")
		}
		jtiDB, err := model.GetJtiById(id) // DBから取得
		if err != nil {
			return nil, "", err
		}
		fmt.Println("jti in claim: " + jti)
		fmt.Println("jti in db: " + jtiDB)
		if jti != jtiDB { // クレームのjtiとusersテーブルのjtiを比較
			return nil, "", errors.New("the jti in the CLAIMS does not match the jti in the user's DB")
		}
		fmt.Println("jtiが一致")
		// expを検証
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, "", errors.New("exp could not be obtained from the token")
		}
		expTT := time.Unix(int64(exp), 0) // Unix 時刻を日時に変換
		timeNow := time.Now()             // 現在時刻を取得
		if timeNow.After(expTT) {         // エラーになるパターン  // 現在時刻timeNowが期限expTTより後ならエラーなのでtrueを出力
			return nil, "", err
		}
	}

	// 正常に終われば解析されたトークンとidを渡す。
	return token, id, nil
}
