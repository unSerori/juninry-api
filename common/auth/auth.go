package auth

import (
	"errors"
	"fmt"
	"juninry-api/common/logging"
	"juninry-api/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ParseTokenの結果用の構造体
type ParseTokenAnalysis struct {
	Token *jwt.Token
	Id    string
}
type Errs struct {
	InputErr    error
	InternalErr error
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
// 成功時に得られる分析結果と複数のエラーが返るので、それぞれ構造体として扱い構造体で返す
// エラーは入力値が不正な場合と処理エラーなどが考えられ、それぞれフィールドとして定義し、呼び出し側では構造体のフィールドを!=nilでチェックしハンドルする
func ParseToken(tokenString string) (ParseTokenAnalysis, Errs) {
	// 返り血用の構造体セット
	var analysis ParseTokenAnalysis
	var errs Errs
	// あらかじめanalysisを宣言したため、analysis, err := ができないことへの対策
	var err error

	// 署名が正しければ、解析用の鍵を使う。(無名関数内で署名方法がHMACであるか確認し、HMACであれば秘密鍵を渡し、jwtトークンを解析する。)
	analysis.Token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // 署名を確認
			logging.ErrorLog(fmt.Sprintf("Unexpected signature method: %v.", token.Header["alg"]), nil)
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil // 署名が正しければJWT_SECRET_KEYをバイト配列にして返す
	})
	if err != nil {
		errs.InternalErr = err
		return analysis, errs
	}

	// 下のクレーム検証処理(:elseスコープ内)で持ち出したい値をあらかじめ宣言しておく。
	//var id string // id
	// 構造体で管理することで不要に！

	// トークン自体が有効か秘密鍵を用いて確認。また、クレーム部分も取得。(トークンの署名が正しいか、有効期限内か、ブラックリストでないか。)
	claims, ok := analysis.Token.Claims.(jwt.MapClaims) // MapClaimsにアサーション
	if !ok || !analysis.Token.Valid {                   // 取得に失敗または検証が失敗
		errs.InputErr = errors.New("invalid authentication token")
		return analysis, errs
	} else { // 有効な場合クレームの各要素を検証
		// idを検証
		idClaims, ok := claims["id"].(string) // goではJSONの数値は少数もカバーしたfloatで解釈される
		if !ok {
			errs.InputErr = errors.New("id could not be obtained from the token")
			return analysis, errs
		}
		analysis.Id = idClaims                           // 調整してスコープ買いに持ち出す。
		if err := model.CfmId(analysis.Id); err != nil { // ユーザーに存在するか。int(id)
			errs.InputErr = err
			return analysis, errs
		}
		// jtiを検証
		jti, ok := claims["jti"].(string)
		if !ok {
			errs.InputErr = errors.New("jti could not be obtained from the token")
			return analysis, errs
		}
		jtiDB, err := model.GetJtiById(analysis.Id) // DBから取得
		if err != nil {
			errs.InputErr = err
			return analysis, errs
		}
		fmt.Println("jti in claim: " + jti)
		fmt.Println("jti in db: " + jtiDB)
		if jti != jtiDB { // クレームのjtiとusersテーブルのjtiを比較
			errs.InputErr = errors.New("the jti in the CLAIMS does not match the jti in the user's DB")
			return analysis, errs
		}
		logging.SuccessLog("JTI matched.")
		// expを検証
		exp, ok := claims["exp"].(float64)
		if !ok {
			errs.InputErr = errors.New("exp could not be obtained from the token")
			return analysis, errs
		}
		expTT := time.Unix(int64(exp), 0) // Unix 時刻を日時に変換
		timeNow := time.Now()             // 現在時刻を取得
		if timeNow.After(expTT) {         // エラーになるパターン  // 現在時刻timeNowが期限expTTより後ならエラーなのでtrueを出力
			errs.InputErr = err
			return analysis, errs
		}
	}

	// 正常に終われば解析されたトークンとidを渡す。
	return analysis, errs
}

// TODO: user auth

// TODO: device auth
