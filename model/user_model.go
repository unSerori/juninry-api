package model

import (
	"juninry-api/logging"
)

// ユーザテーブル  // モデルを構造体で定義
type User struct { // typeで型の定義, structは構造体
	UserUuid    string  `xorm:"varchar(36) pk" json:"userUUID"`                  // ユーザのUUID
	UserName    string  `xorm:"varchar(25) not null" json:"userName"`            // 名前
	UserTypeId  int     `xorm:"not null" json:"userTypeId"`                      // ユーザータイプ
	MailAddress string  `xorm:"varchar(256) not null unique" json:"mailAddress"` // メアド
	Password    string  `xorm:"varchar(60) not null" json:"password"`            // bcrypt化されたパスワード
	JtiUuid     string  `xorm:"varchar(36) unique" json:"jwtUUID"`               // jwtクレームのuuid
	OuchiUuid   *string `xorm:"varchar(36) default NULL" json:"ouchiUUID"`       // 所属するおうちのUUID
}

// テーブル名
func (User) TableName() string {
	return "users"
}

// FK制約の追加
func InitUserFK() error {
	// UserTypeId
	_, err := db.Exec("ALTER TABLE users ADD FOREIGN KEY (user_type_id) REFERENCES user_types(user_type_id) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// OuchiUuid
	_, err = db.Exec("ALTER TABLE users ADD FOREIGN KEY (ouchi_uuid) REFERENCES ouchies(ouchi_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateUserTestData() {
	user4 := &User{
		UserUuid:    "3cac1684-c1e0-47ae-92fd-6d7959759224",
		UserName:    "test pupil",
		UserTypeId:  2,
		MailAddress: "test-pupil@gmail.com",
		Password:    "$2a$10$8hJGyU235UMV8NjkozB7aeHtgxh39wg/ocuRXW9jN2JDdO/MRz.fW", // C@tp
		JtiUuid:     "14dea318-8581-4cab-b233-995ce8e1a948",
	}
	db.Insert(user4)
	user5 := &User{
		UserUuid:    "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		UserName:    "test teacher",
		UserTypeId:  1,
		MailAddress: "test-teacher@gmail.com",
		Password:    "$2a$10$Ig/s1wsrXBuZ7qvjudr4CeQFhqJTLQpoAAp1LrBNh5jX9VZZxa3R6", // C@tt
		JtiUuid:     "42c28ac4-0ba4-4f81-8813-814dc92e2f40",
	}
	db.Insert(user5)
	user6 := &User{
		UserUuid:    "868c0804-cf1b-43e2-abef-08f7ef58fcd0",
		UserName:    "test parent",
		UserTypeId:  3,
		MailAddress: "test-parent@gmail.com",
		Password:    "$2a$10$8hJGyU235UMV8NjkozB7aeHtgxh39wg/ocuRXW9jN2JDdO/MRz.fW", // C@tp
		JtiUuid:     "0553853f-cbcf-49e2-81d6-a4c7e4b1b470",
	}
	db.Insert(user6)
}

// 新規ユーザ登録
// 新しい構造体をレコードとして受け取り、usersテーブルにinsertし、成功した列数とerrorを返す
func CreateUser(record User) (int64, error) {
	affected, err := db.Insert(record)
	return affected, err
}

// jtiを保存更新
// user_uuidと更新用のjtiを受け取り、jti_uuid列を更新し、エラーがあれば返す
func SaveJti(userUuid string, jti string) error {
	// jtiを更新
	_, err := db.Cols("jti_uuid").Where("user_uuid = ?", userUuid).Update(&User{JtiUuid: jti}) // update時にはWhereと列制限を使うとよい  // Omit<->Cols
	if err != nil {
		return err
	}

	return nil
}

// ユーザー情報の取得
// user_uuidを受け取り、usersテーブルから該当するレコードを取得し、構造体にマッピングして返す
func GetUser(userUuid string) (User, error) {
	var user User
	_, err := db.Where("user_uuid = ?", userUuid).Get(&user)
	return user, err
}

// idが存在するか確かめる
func CfmId(userUuid string) error {
	var user User // 取得したデータをマッピングする構造体
	_, err := db.Where("user_uuid = ?", userUuid).Get(&user)
	if err != nil {
		return err
	}
	return nil // エラーなければnilが返る
}

// ユーザのuuidからjtiを取得
func GetJtiById(userUuid string) (string, error) {
	var user User // 取得したデータをマッピングする構造体

	// 該当ユーザの行を取得
	_, err := db.Where("user_uuid = ?", userUuid).Get(&user)
	if err != nil {
		return "", err
	}

	return user.JtiUuid, nil
}

// メアドからユーザーが存在するか確認
func CheckUserExists(mail string) (error, bool) {
	var user User // 取得したデータをマッピングする構造体

	isFound, err := db.Where("mail_address = ?", mail).Get(&user)
	if err != nil {
		logging.ErrorLog("Error when searching for a user from a mail address.", err)
		return err, isFound
	}
	if !isFound {
		logging.ErrorLog("Could not find any users from the e-mail address.", err)
		return nil, isFound
	}

	return nil, true
}

// メアドからパスワードを取得
func GetPassByMail(mail string) (string, error, bool) {
	var user User // 取得したデータをマッピングする構造体

	isFound, err := db.Select(
		"password", // パスワードをとる
	).Where("mail_address = ?", mail).Get(&user) // Select(必要な列).Where(会社番号が引数の値).Find(User構造体の形で取得)
	if err != nil {
		return "", err, isFound
	}
	if !isFound { // 見つからなかった
		return "", nil, false
	}

	return user.Password, nil, true // ユーザースライスを返す。
}

// メアドからuuidを取得
func GetIdByMail(mail string) (string, error, bool) {
	var user User // 取得したデータをマッピングする構造体

	isFound, err := db.Select("user_uuid").Where("mail_address = ?", mail).Get(&user)
	if err != nil {
		return "", err, isFound
	}
	if !isFound { // 見つからなかった
		return "", nil, false
	}

	return user.UserUuid, nil, true
}

// アカウントタイプが教師かどうか判定して真偽値を返す
func IsTeacher(userUuid string) (bool, error) {
	var user User // 取得したデータをマッピングする構造体
	// TODO: 教員のみに制限する
	// 該当ユーザの行を取得
	isTeacher, err := db.Where("user_uuid = ? and user_type_id = 1", userUuid).Exist(&user)
	if err != nil {
		return false, err // エラーが出てるのにfalse返すのきしょいかも
	}

	return isTeacher, nil
}

// アカウントタイプが親かどうか判定して真偽値を返す
func IsPatron(userUuid string) (bool, error) {
	var user User // 取得したデータをマッピングする構造体
	// 該当ユーザの行を取得
	isParent, err := db.Where("user_uuid = ? and user_type_id = 3", userUuid).Exist(&user)
	if err != nil {
		return false, err // エラーが出てるのにfalse返すのきしょいかも
	}

	return isParent, nil
}
