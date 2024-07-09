package service

import (
	"crypto/rand"
	"fmt"
	"juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"math/big"
	"time"

	"github.com/go-sql-driver/mysql"
)

type OuchiService struct{}

// 招待コード生成部分
// クラス内でしか呼び出されない
func (s *OuchiService) generateOuchiInviteCode(bOuchi model.Ouchi) (model.Ouchi, error) {
	// 有効な招待コードが無ければ新しい招待コードを作る
	// 有効期限を1週間後に設定
	validUntil := time.Now().AddDate(0, 0, 7)
	bOuchi.ValidUntil = validUntil // バインド

	// 10回エラー吐いたら終わりでええやろ。。。
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		// 招待コードを作る
		inviteCode, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil { // 乱数生成でエラーが出たら泣く
			continue
		}
		// 4桁文字列にキャストしてバインド
		bOuchi.InviteCode = fmt.Sprintf("%04d", inviteCode.Int64())

		// クラステーブルに追加
		_, err = model.UpdateOuchiInviteCode(bOuchi)
		if err != nil {
			// 本処理時のエラーごとに処理(:DBエラー)
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反の場合、招待コードから再生成
				continue
			default: // 一意性制約違反じゃなかったらびっくり
				return model.Ouchi{}, err
			}
		}

		// 招待コードを作成できたクラスを返す
		return bOuchi, nil
	}

	// 試行回数10回以上で失敗したらエラーを返す
	// これ10回連続衝突する可能性そこそこあるよね〜
	// TODO: 改善の余地あり
	logging.ErrorLog("Maximum number of attempts reached", nil)
	return model.Ouchi{}, common.NewErr(common.ErrTypeMaxAttemptsReached)
}

