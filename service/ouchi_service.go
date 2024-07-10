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
	"github.com/google/uuid"
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
		// 6桁文字列にキャストしてバインド
		bOuchi.InviteCode = fmt.Sprintf("%06d", inviteCode.Int64())

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

func (s *OuchiService) PermissionCheckedOuchiCreation(userUuid string, bOuchi model.Ouchi) (model.Ouchi, error) {

	// おうち作成権限を持っているか確認
	isPatron, err := model.IsPatron(userUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	fmt.Println(isPatron)
	if !isPatron { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return model.Ouchi{}, common.NewErr(common.ErrTypePermissionDenied)
	}

	// ユーザがすでにおうちに所属していないかの確認
	user, err := model.GetUser(userUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	// null確認
	if user.OuchiUuid != nil {
		logging.ErrorLog("You are already assigned to an Ouchi", nil)
		return model.Ouchi{}, err
	}

	// おうち作成
	// おうちUUIDの生成
	ouchiUuid, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {                    // 空の構造体とエラー
		return model.Ouchi{}, err
	}

	bOuchi.OuchiUuid = ouchiUuid.String() // バインド

	// おうち作成
	_, err = model.CreateOuchi(bOuchi)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err // uuidの作成がおかしくなければ問題ないけど、登録結果が0件で正常終了することなんかあるか？
	}

	// 招待コード入ったクラスもらえます！
	ouchi, err := s.generateOuchiInviteCode(bOuchi)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}

	// 保護者にouchiUuidを付与
	_, err = model.AssignOuchi(userUuid, bOuchi.OuchiUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}

	fmt.Println("確認")
	fmt.Println(ouchi)

	//エラーが出なかった場合、コミットして作成したクラスを返す
	return ouchi, nil
}

func (s *OuchiService) PermissionCheckedRefreshOuchiInviteCode(userUuid string, ouchiUuid string) (model.Ouchi, error) {

	// クラス作成権限を持っているか確認
	isPatron, err := model.IsPatron(userUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err // トークンあるのにユーザーがいないことはあり得ないのでないと思うが、、、？
	}
	if !isPatron { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return model.Ouchi{}, common.NewErr(common.ErrTypePermissionDenied)
	}
	// クラスUUIDが存在するかどうか
	targetouchi, err := model.GetOuchi(ouchiUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	if targetouchi.OuchiUuid == "" { // そんなクラス存在しない場合
		return model.Ouchi{}, common.NewErr(common.ErrTypeNoResourceExist)
	}

	// 招待コード入ったクラスもらえます！
	ouchi, err := s.generateOuchiInviteCode(targetouchi)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	//エラーが出なかった場合、コミットして作成したクラスを返す
	return ouchi, nil
}
