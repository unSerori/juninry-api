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
// おうち内でしか呼び出されない(頭文字が小文字の場合プライベート)
func (s *OuchiService) generateOuchiInviteCode(bOuchi model.Ouchi) (model.Ouchi, error) {
	// 有効な招待コードが無ければ新しい招待コードを作る
	// 有効期限を1週間後に設定
	validUntil := time.Now().AddDate(0, 0, 1)
	bOuchi.ValidUntil = validUntil // バインド

	// 10回エラー吐いたら終わり
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		// 招待コードを作る
		inviteCode, err := rand.Int(rand.Reader, big.NewInt(1000000))
		if err != nil {
			continue // breakみたいなもん？
		}
		// 6桁文字列にキャストしてバインド
		bOuchi.InviteCode = fmt.Sprintf("%06d", inviteCode.Int64())

		// おうちテーブルに追加
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
	logging.ErrorLog("Maximum number of attempts reached", nil)
	return model.Ouchi{}, common.NewErr(common.ErrTypeMaxAttemptsReached)
}

// おうち作成
func (s *OuchiService) PermissionCheckedOuchiCreation(userUuid string, bOuchi model.Ouchi) (model.Ouchi, error) {

	// おうち作成権限を持っているか確認(親)
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

	// おうちUUIDの生成
	ouchiUuid, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {                    // 空の構造体とエラー
		return model.Ouchi{}, err
	}

	bOuchi.OuchiUuid = ouchiUuid.String() // バインド

	// おうち作成
	_, err = model.CreateOuchi(bOuchi)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
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

	//エラーが出なかった場合、コミットして作成したおうちを返す
	return ouchi, nil
}

// 招待コード更新処理
func (s *OuchiService) PermissionCheckedRefreshOuchiInviteCode(userUuid string, ouchiUuid string) (model.Ouchi, error) {

	// おうち作成権限を持っているか確認
	isPatron, err := model.IsPatron(userUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	if !isPatron { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return model.Ouchi{}, common.NewErr(common.ErrTypePermissionDenied)
	}
	// おうちUUIDが存在するかどうか
	targetouchi, err := model.GetOuchi(ouchiUuid)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	if targetouchi.OuchiUuid == "" { // そんなおうち存在しない場合弾く
		return model.Ouchi{}, common.NewErr(common.ErrTypeNoResourceExist)
	}

	// 招待コード入ったクラスもらえます！
	ouchi, err := s.generateOuchiInviteCode(targetouchi)
	if err != nil { // エラーハンドル
		return model.Ouchi{}, err
	}
	//エラーが出なかった場合、コミットして作成したおうちを返す
	return ouchi, nil
}

// おうちに所属処理
func (s *OuchiService) PermissionCheckedJoinOuchi(userUuid string, inviteCode string) (string, error) {

	// おうちに所属できるuserTypeか確認
	isJunior, err := model.IsJunior(userUuid)
	if err != nil { // エラーハンドル
		return "", nil
	}
	if !isJunior { // 先生も親も両方弾く
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return "", common.NewErr(common.ErrTypePermissionDenied)
	}

	// ouchiUuidが存在するか
	targetOuchi, err := model.GetOuchiInviteCode(inviteCode)
	if err != nil { // エラーハンドル
		return "", err
	}
	if targetOuchi.OuchiUuid == "" { // おうちが存在しない場合
		return "", common.NewErr(common.ErrTypeNoResourceExist)
	}

	// ユーザにouchiUuidを付与
	_, err = model.AssignOuchi(userUuid, targetOuchi.OuchiUuid)
	if err != nil { // エラーハンドル
		return "", err
	}

	//　所属したおうちを返す
	ouchi, err := model.GetOuchi(targetOuchi.OuchiUuid)
	if err != nil {
		return "", nil
	}

	// エラーがない場合、おうち名返還
	return ouchi.OuchiName, nil
}
