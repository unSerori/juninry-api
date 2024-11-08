package service

import (
	"fmt"
	"juninry-api/common/auth"
	"juninry-api/common/custom"
	"juninry-api/common/logging"
	"juninry-api/common/security"
	"juninry-api/model"
	"time"

	"github.com/google/uuid"
)

type UserService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

func (s *UserService) RegisterUser(bUser model.User) (string, error) {
	// user_uuidを生成
	userId, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {
		return "", err
	}
	bUser.UserUuid = userId.String() // 設定

	// パスワードをハッシュ化
	hashed, err := security.HashingByEncrypt(bUser.Password)
	if err != nil {
		return "", custom.NewErr(custom.ErrTypeHashingPassFailed, custom.WithMsg(err.Error())) // controller側でエラーを判別するために、ENUMで管理されたエラーを返す
	}
	// ハッシュ(:バイト配列)化されたパスワードを文字列にsh知恵構造体に戻す
	bUser.Password = string(hashed)

	// 構造体をレコード登録処理に投げる
	err = model.CreateUser(bUser) // 第一返り血は登録成功したレコード数
	if err != nil {               // エラーハンドル
		return "", err
	}

	// 登録が成功したらトークンを生成する
	token, err := auth.GenerateToken(bUser.UserUuid) // トークンを取得
	if err != nil {
		return "", custom.NewErr(custom.ErrTypeGenTokenFailed, custom.WithMsg(err.Error()))
	}

	// TODO: userTypeをそのまま取ってこれるのに、生徒チェックする必要ある？こっちの書き方の方が見たとき分かりやすよね
	isJunior, err := model.IsJunior(bUser.UserUuid)
	if err != nil {
		return "", err
	}
	// 生徒だったらスタンプカードを作成する＋メインにゃいおっとを生成
	if isJunior {
		// insert用のレコード作成(Stamp)
		var bStamp model.Stamp

		// 現在の日時を取得
		now := time.Now()
		// 2日前の日付を取得
		yesterday := now.AddDate(0, 0, -2)
		fmt.Printf("yesterday: %v\n", yesterday)
		// 日付部分だけを取得（時間は00:00:00）
		dateOnly := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

		//　userUuid,quentity,loginTimeを追加
		bStamp.UserUuid = bUser.UserUuid
		bStamp.LastLoginTime = dateOnly

		// スタンプカード作る
		_, err := model.CreateStampCard(bStamp)
		if err != nil {
			return "", err
		}

		fmt.Println("junia's stamp cards are ready!!!!")

		// ニャリオット所持テーブルに値入れる
		_, err = model.CreateNyariotInventory(model.NyariotInventory{
			UserUuid:     bStamp.UserUuid,
			NyariotUuid:  "c0768960-eb5f-4a60-8327-4171fd4b8a46",
			ConvexNumber: 1,
		})
		if err != nil {
			return "", err
		}
		fmt.Print("Nyariot possession registration completed!!!!")

		// ハングリーテーブル作る
		_, err = model.CreateHungryStatus(model.HungryStatus{
			UserUuid:      bStamp.UserUuid,
			SatityDegrees: 100,
			LastGohanTime: time.Now(),
			NyariotUuid:   "c0768960-eb5f-4a60-8327-4171fd4b8a46",
		})
		if err != nil {
			return "", err
		}

		fmt.Print("completed the Nyariot setup!!!!")
	}

	return token, nil
}

// ユーザ情報取得
func (s *UserService) GetUser(useruuid string) (model.User, error) {

	// useridでユーザ情報を取得
	user, err := model.GetUser(useruuid)
	if err != nil {
		return user, err
	}

	return user, err
}

// ログイン
func (s *UserService) LoginUser(bUser model.User) (string, error) {
	// ユーザーの存在確認
	err, isFound := model.CheckUserExists(bUser.MailAddress)
	if err != nil {
		logging.ErrorLog("No corresponding user exists.", err)
		return "", err
	}
	if !isFound {
		logging.ErrorLog("Could not find the relevant ID.", nil)
		return "", custom.NewErr(custom.ErrTypeNoResourceExist)
	}

	// 登録済みのパスワードを取得し、
	pass, err, isFound := model.GetPassByMail(bUser.MailAddress)
	if err != nil {
		logging.ErrorLog("Failed to retrieve password.", err)
		return "", err
	}
	if !isFound {
		logging.ErrorLog("Could not find the relevant ID.", nil)
		return "", custom.NewErr(custom.ErrTypeNoResourceExist)
	}
	// 比較する
	err = security.CompareHashAndStr([]byte(pass), bUser.Password)
	if err != nil {
		logging.ErrorLog("Password does not match.", err)
		return "", custom.NewErr(custom.ErrTypePassMismatch, custom.WithMsg(err.Error()))
	}

	// トークンを生成しなおす
	id, err, isFound := model.GetIdByMail(bUser.MailAddress) // user_uuidを取得し、
	if err != nil {
		logging.ErrorLog("Failure to obtain id.", err)
		return "", err
	}
	if !isFound { // idが見つからなかった
		logging.ErrorLog("Could not find the relevant ID.", err)
		return "", custom.NewErr(custom.ErrTypeNoResourceExist)
	}
	token, err := auth.GenerateToken(id) // user_uuidをもとにトークンを生成
	if err != nil {
		logging.ErrorLog("Failed to generate token.", err)
		return "", err
	}

	return token, nil
}

// ポイントをアップデート
