package service

import (
	"juninry-api/auth"
	"juninry-api/common"
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
		return "", common.NewErr(common.ErrTypeHashingPassFailed, common.WithMsg(err.Error())) // controller側でエラーを判別するために、ENUMで管理されたエラーを返す
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
		return "", common.NewErr(common.ErrTypeGenTokenFailed, common.WithMsg(err.Error()))
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
		return "", common.NewErr(common.ErrTypeNoResourceExist)
	}

	// 登録済みのパスワードを取得し、
	pass, err, isFound := model.GetPassByMail(bUser.MailAddress)
	if err != nil {
		logging.ErrorLog("Failed to retrieve password.", err)
		return "", err
	}
	if !isFound {
		logging.ErrorLog("Could not find the relevant ID.", nil)
		return "", common.NewErr(common.ErrTypeNoResourceExist)
	}
	// 比較する
	err = security.CompareHashAndStr([]byte(pass), bUser.Password)
	if err != nil {
		logging.ErrorLog("Password does not match.", err)
		return "", common.NewErr(common.ErrTypePassMismatch, common.WithMsg(err.Error()))
	}

	// トークンを生成しなおす
	id, err, isFound := model.GetIdByMail(bUser.MailAddress) // user_uuidを取得し、
	if err != nil {
		logging.ErrorLog("Failure to obtain id.", err)
		return "", err
	}
	if !isFound { // idが見つからなかった
		logging.ErrorLog("Could not find the relevant ID.", err)
		return "", common.NewErr(common.ErrTypeNoResourceExist)
	}
	token, err := auth.GenerateToken(id) // user_uuidをもとにトークンを生成
	if err != nil {
		logging.ErrorLog("Failed to generate token.", err)
		return "", err
	}

	return token, nil
}
