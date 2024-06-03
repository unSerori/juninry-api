package model

import "fmt"

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
	User4 := &User{
		UserUuid:    "3cac1684-c1e0-47ae-92fd-6d7959759224",
		UserName:    "test pupil",
		UserTypeId:  2,
		MailAddress: "test-pupil@gmail.com",
		Password:    "$2a$10$8hJGyU235UMV8NjkozB7aeHtgxh39wg/ocuRXW9jN2JDdO/MRz.fW", // C@tp
		JtiUuid:     "14dea318-8581-4cab-b233-995ce8e1a948",
	}
	db.Insert(User4)
	User5 := &User{
		UserUuid:    "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		UserName:    "test teacher",
		UserTypeId:  1,
		MailAddress: "test-teacher@gmail.com",
		Password:    "$2a$10$Ig/s1wsrXBuZ7qvjudr4CeQFhqJTLQpoAAp1LrBNh5jX9VZZxa3R6", // C@tt
		JtiUuid:     "42c28ac4-0ba4-4f81-8813-814dc92e2f40",
	}
	_, err := db.Insert(User5)
	if err != nil {
		fmt.Println(err)
	}

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
