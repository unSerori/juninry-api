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

type ClassService struct{}

func (s *ClassService) PermissionCheckedClassCreation(userUuid string, bClass model.Class) (model.Class, error) {

	// クラス作成権限を持っているか確認
	isTeacher, err := model.CheckIsTeacher(userUuid)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return model.Class{}, common.NewErr(common.ErrTypePermissionDenied)
	}

	// クラス作成
	// クラスUUIDの生成
	classUuid, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {                    // 空の構造体とエラー
		return model.Class{}, err
	}

	bClass.ClassUuid = classUuid.String() // バインド

	// クラス作成
	_, err = model.CreateClass(bClass)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}

	// 招待コード生成
	bClass, err = s.generateInviteCode(bClass)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}

	// 所属構造体にユーザーIDとクラスIDをバインド
	classMemberships := model.ClassMembership{
		ClassUuid: bClass.ClassUuid,
		UserUuid:  userUuid,
	}

	// 教員をクラスに所属させる
	_, err = model.JoinClass(classMemberships)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}

	//エラーが出なかった場合、コミットして作成したクラスを返す
	return bClass, nil

}

func (s *ClassService) generateInviteCode(bClass model.Class) (model.Class, error) {

	// 有効な招待コードが無ければ新しい招待コードを作る
	// 有効期限を1週間後に設定
	validUntil := time.Now().AddDate(0, 0, 7)
	bClass.ValidUntil = validUntil // バインド

	// 10回エラー吐いたら終わりでええやろ。。。
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		// 招待コードを作る
		inviteCode, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil { // 乱数生成でエラーが出たら泣く
			continue
		}
		// 4桁文字列にキャストしてバインド
		bClass.InviteCode = fmt.Sprintf("%04d", inviteCode.Int64())

		// クラステーブルに追加
		_, err = model.UpdateInviteCode(bClass)
		if err != nil {
			// 本処理時のエラーごとに処理(:DBエラー)
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反の場合、招待コードから再生成
				continue
			default: // 一意性制約違反じゃなかったらびっくり
				return model.Class{}, err
			}
		}
		// クラス返して
		return bClass, nil
	}

	// 試行回数10回以上で失敗したらエラーを返す
	// これ10回連続衝突する可能性そこそこあるよね〜
	// TODO: 改善の余地あり
	logging.ErrorLog("Maximum number of attempts reached", nil)
	return model.Class{}, common.NewErr(common.ErrTypeMaxAttemptsReached)
}
