package service

import (
	"errors"
	"juninry-api/common/logging"
	"juninry-api/model"
	"juninry-api/utility/auth"
	"juninry-api/utility/custom"
	"juninry-api/utility/security"

	"github.com/go-sql-driver/mysql"
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
	_, err = model.CreateUser(bUser) // 第一返り血は登録成功したレコード数
	if err != nil {                  // エラーハンドル
		// XormのORMエラーを仕分ける
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // errをmysqlErrにアサーション出来たらtrue
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反
				return "", custom.NewErr(custom.ErrTypeUniqueConstraintViolation)
			default: // ORMエラーの仕分けにぬけがある可能性がある
				return "", custom.NewErr(custom.ErrTypeOtherErrorsInTheORM)
			}
		}
		// 通常の処理エラー
		return "", err
	}

	// 登録が成功したらトークンを生成する
	token, err := auth.GenerateToken(bUser.UserUuid) // トークンを取得
	if err != nil {
		return "", custom.NewErr(custom.ErrTypeGenTokenFailed, custom.WithMsg(err.Error()))
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
