package model

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

// テーブル名
func (User) TableName() string {
	return "users"
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
