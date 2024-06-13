package service

import (
	"juninry-api/auth"
	commons "juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/security"

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
		return "", commons.NewErr(commons.ErrTypeHashingPassFailed, err.Error()) // controller側でエラーを判別するために、ENUMで管理されたエラーを返す
	}
	// ハッシュ(:バイト配列)化されたパスワードを文字列にsh知恵構造体に戻す
	bUser.Password = string(hashed)

	// 構造体をレコード登録処理に投げる
	_, err = model.CreateUser(bUser) // 第一返り血は登録成功したレコード数
	if err != nil {
		return "", err
	}

	// 登録が成功したらトークンを生成する
	token, err := auth.GenerateToken(bUser.UserUuid) // トークンを取得
	if err != nil {
		return "", commons.NewErr(commons.ErrTypeGenTokenFailed, err.Error())
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
	err := model.CheckUserExists(bUser.MailAddress)
	if err != nil {
		logging.ErrorLog("No corresponding user exists.", nil)
		return "", err
	}

	// 登録済みのパスワードを取得し、
	pass, err := model.GetPassByMail(bUser.MailAddress)
	if err != nil {
		logging.ErrorLog("Failed to retrieve password.", nil)
		return "", err
	}
	// 比較する
	err = security.CompareHashAndStr([]byte(pass), bUser.Password)
	if err != nil {
		logging.ErrorLog("Password does not match.", nil)
		return "", err
	}

	// トークンを生成しなおす
	id, err := model.GetIdByMail(bUser.MailAddress) // user_uuidを取得し、
	if err != nil {
		return "", err
	}
	token, err := auth.GenerateToken(id) // user_uuidをもとにトークンを生成
	if err != nil {
		return "", err
	}

	return token, nil
}
