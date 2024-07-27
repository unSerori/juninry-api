package security

import (
	"golang.org/x/crypto/bcrypt"
)

// ハッシュ化
func HashingByEncrypt(str string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost) // 第一引数: ハッシュ化する平分の文字列をバイト列にして渡す  第二引数: ストレッチング回数(Defaultは10回<-2^N)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

// ハッシュ化されている値(:バイト配列)と平文を比較
func CompareHashAndStr(hashed []byte, str string) error {
	return bcrypt.CompareHashAndPassword(hashed, []byte(str)) // バイト列で比較、返り血はエラー
}
