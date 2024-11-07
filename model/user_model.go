package model

import (
	"errors"
	"juninry-api/common/logging"
	"juninry-api/utility/custom"

	"github.com/go-sql-driver/mysql"
)

// ユーザテーブル  // モデルを構造体で定義
type User struct { // typeで型の定義, structは構造体
	UserUuid    string  `xorm:"varchar(36) pk" json:"userUUID"`                  // ユーザのUUID
	UserName    string  `xorm:"varchar(25) not null" json:"userName"`            // 名前
	UserTypeId  int     `xorm:"not null" json:"userTypeId"`                      // ユーザータイプ	1:教師, 2:児童, 3:保護者
	MailAddress string  `xorm:"varchar(256) not null unique" json:"mailAddress"` // メアド
	GenderId    int     `xorm:"not null" json:"genderId"`                        // 性別 1:男性, 2:女性, 3:その他
	Password    string  `xorm:"varchar(60) not null" json:"password"`            // bcrypt化されたパスワード
	JtiUuid     string  `xorm:"varchar(36) unique" json:"jwtUUID"`               // jwtクレームのuuid
	OuchiUuid   *string `xorm:"varchar(36) default NULL" json:"ouchiUUID"`       // 所属するおうちのUUID
	OuchiPoint  int     `xorm:"default 0" json:"ouchiPoint"`                     // おうちのポイント
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

	str := "2e17a448-985b-421d-9b9f-62e5a4f28c49"
	strPtr := &str

	// ポインタ型の変数に割り当て
	var ouchiUUID *string = strPtr

	user3 := &User{
		UserUuid:    "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		UserName:    "test teacher",
		UserTypeId:  1,
		MailAddress: "test-teacher@gmail.com",
		Password:    "$2a$10$Ig/s1wsrXBuZ7qvjudr4CeQFhqJTLQpoAAp1LrBNh5jX9VZZxa3R6", // C@tt
		JtiUuid:     "42c28ac4-0ba4-4f81-8813-814dc92e2f40",
	}
	db.Insert(user3)

	user4 := &User{
		UserUuid:    "3cac1684-c1e0-47ae-92fd-6d7959759224",
		UserName:    "test pupil",
		UserTypeId:  2,
		GenderId:    1,
		MailAddress: "test-pupil@gmail.com",
		Password:    "$2a$10$8hJGyU235UMV8NjkozB7aeHtgxh39wg/ocuRXW9jN2JDdO/MRz.fW", // C@tp
		JtiUuid:     "14dea318-8581-4cab-b233-995ce8e1a948",
		OuchiUuid:   ouchiUUID,
	}
	db.Insert(user4)

	user5 := &User{
		UserUuid:    "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		UserName:    "test teacher",
		UserTypeId:  1,
		GenderId:    2,
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
		OuchiUuid:   ouchiUUID,
	}
	db.Insert(user6)

	user7 := &User{
		UserUuid:    "cd09ac2f-4278-4fb0-a8bc-df7c2d9ef1fc",
		UserName:    "test pupil2go",
		UserTypeId:  2,
		GenderId:    1,
		MailAddress: "test-pupil2go@gmail.com",
		Password:    "$2a$10$8hJGyU235UMV8NjkozB7aeHtgxh39wg/ocuRXW9jN2JDdO/MRz.fW", // C@tp
		JtiUuid:     "b8595062-c70a-48ee-be8f-dce768d49675",
		OuchiUuid:   ouchiUUID,
	}
	db.Insert(user7)
}

// 新規ユーザ登録
// 新しい構造体をレコードとして受け取り、usersテーブルにinsertし、成功した列数とerrorを返す
func CreateUser(record User) error {
	affected, err := db.Insert(record)
	if err != nil { //エラーハンドル
		// XormのORMエラーを仕分ける
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // errをmysqlErrにアサーション出来たらtrue
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反
				return custom.NewErr(custom.ErrTypeUniqueConstraintViolation)
			default: // ORMエラーの仕分けにぬけがある可能性がある
				return custom.NewErr(custom.ErrTypeOtherErrorsInTheORM)
			}
		}
		// 通常の処理エラー
		return err
	}
	if affected == 0 {
		return custom.NewErr(custom.ErrTypeZeroEffectCUD)
	}

	return nil
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

// 複数件のユーザー情報の取得
func GetUsers(userUuid []string) ([]User, error) {
	var users []User
	err := db.In("user_uuid", userUuid).And("user_type_id = 2").Find(&users)
	return users, err
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

// uuidから名前を取得
func GetNameById(userId string) (string, error) {
	var user User // 取得したデータをマッピングする構造体

	// 該当ユーザの行を取得
	_, err := db.Where("user_uuid = ?", userId).Get(&user)
	if err != nil {
		return "", err
	}

	return user.UserName, nil
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

// IDからアカウントタイプを返す
func GetUserTypeId(userId string) (int, error) {
	var user User // 取得したデータをマッピングする構造体

	// 該当ユーザーを列ごと取得
	isFound, err := db.Where("user_uuid = ?", userId).Get(&user)
	if err != nil { //エラーハンドル
		return 0, err
	}
	if !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`!isFound`で、影響を与えるSQL文の時は`affected == 0`でハンドリング
		return 0, custom.NewErr(custom.ErrTypeNoFoundR)
	}

	return user.UserTypeId, nil // teacher, junior, patron
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
	isPatron, err := db.Where("user_uuid = ? and user_type_id = 3", userUuid).Exist(&user)
	if err != nil {
		return false, err // エラーが出てるのにfalse返すのきしょいかも
	}

	return isPatron, nil
}

// アカウントタイプがおこさまかどうか判定して真偽値を返す
func IsJunior(userUuid string) (bool, error) {
	var user User // 取得したデータをマッピングする構造体
	// TODO: ガキのみに制限する
	// 該当ユーザの行を取得
	isJunior, err := db.Where("user_uuid = ? and user_type_id = 2", userUuid).Exist(&user)
	if err != nil {
		return false, err // エラーが出てるのにfalse返すのきしょいかも
	}

	return isJunior, nil
}

// 子供のUUIDを取得
func GetChildrenUuids(OuchiUuid string) ([]string, error) {

	// 結果格納用変数
	var userUuids []string

	err := db.Table("users").Where("ouchi_uuid = ? and user_type_id = 2", OuchiUuid).Select("user_uuid").Find(&userUuids)
	if err != nil {
		return nil, err
	}

	return userUuids, nil
}

// ユーザにouchiUuidを付与
func AssignOuchi(userUuid string, ouchiUuid string) (int64, error) {
	// ouchiUuidフィールドにポインタを指定
	user := User{OuchiUuid: &ouchiUuid}
	// 付与処理（更新処理）
	affected, err := db.Where("user_uuid = ?", userUuid).Update(&user)
	return affected, err
}

// ouchiUuidとclassUuidからおこさまを取得
func GetJunior(ouchiUuid string) (User, error) {
	//junior
	var junior User
	_, err := db.Where("ouchi_uuid = ? and user_type_id = 2", ouchiUuid).Get(&junior)
	return junior, err
}

// helpをもとにポイントを加算
func IncrementUpdatePoint(userUuid string, helpUUID string) (*int, error) {

	// 現在のポイントを取得
	user, err := GetUser(userUuid)
	if err != nil {
		return nil, err
	}
	// おてつだいを取得
	help, err := GetHelp(helpUUID)
	if err != nil {
		return nil, err
	}

	incrementedPoint := user.OuchiPoint + help.RewardPoint
	// ポイントを更新
	_, err = db.Cols("ouchi_point").Where("user_uuid = ?", userUuid).Update(&User{OuchiPoint: incrementedPoint})
	if err != nil {
		return nil, err
	}
	ouchiPoint := &incrementedPoint
	return ouchiPoint, err
}

// rewardをもとにポイントを減算
func DecrementUpdatePoint(userUuid string, rewardUUID string) (int, error) {
	// 現在のポイントを取得
	user, err := GetUser(userUuid)
	if err != nil {
		return 0, err
	}
	// ごほうびを取得
	reward, err := GetReward(rewardUUID)
	if err != nil {
		return 0, err
	}
	decrementedPoint := user.OuchiPoint - reward.RewardPoint
	// ポイントを更新
	_, err = db.Cols("ouchi_point").Where("user_uuid = ?", userUuid).Update(&User{OuchiPoint: decrementedPoint})
	if err != nil {
		return 0, err
	}
	return decrementedPoint, err
}

// 同じouchiUuidの人を取得
func GetUserByOuchiUuid(ouchiUuid string) ([]User, error) {
	//ユーザを取得
	var users []User
	err := db.In("ouchi_uuid", ouchiUuid).OrderBy("user_type_id DESC").Find(&users)
	return users, err
}

// IDからouchiUuidを取得
func GetOuchiUuidById(userId string) (string, error) {
	var user User // 取得したデータをマッピングする構造体

	// 該当ユーザの行を取得
	_, err := db.Where("user_uuid = ?", userId).Get(&user)
	if err != nil {
		return "", err
	}

	return *user.OuchiUuid, nil
}

// おうちに所属する保護者を取得
func GetPatronByOuchiUuid(ouchiId string) (User, error) {
	var patron User // 取得したデータをマッピングする構造体
	isFound, err := db.Where("ouchi_uuid = ?", ouchiId).Get(&patron)
	if err != nil {
		return User{}, err
	}
	if !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`!isFound`で、影響を与えるSQL文の時は`affected == 0`でハンドリング
		return User{}, custom.NewErr(custom.ErrTypeNoFoundR)
	}
	return patron, nil // エラーなければnilが返る
}

// おうちに所属するおこさまたちを取得
func GetJuniorsByOuchiUuid(ouchiId string) ([]User, error) {
	var juniors []User // 取得したデータをマッピングする構造体
	err := db.Where("ouchi_uuid = ? and user_type_id = 2", ouchiId).Find(&juniors)
	if err != nil {
		return []User{}, nil
	}
	return juniors, nil // エラーなければnilが返る
}

// 消費可能な最大ポイントを取得
func GetOuchiPointByUserUuid(userUuid string) (int, error) {

	// 結果格納用変数
	var user User

	// 消費可能な最大ポイントを取得
	_, err := db.Where("user_uuid = ?", userUuid).Get(&user)
	if err != nil {
		return 0, err
	}

	return user.OuchiPoint, nil
}

// ポイントを更新
func UpdateOuchiPoint(userUuid string, point int) error {
	_, err := db.Cols("ouchi_point").Where("user_uuid = ?", userUuid).Update(&User{OuchiPoint: point})
	if err != nil {
		return err
	}
	return nil
}
