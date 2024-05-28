package model

// ユーザテーブル  // モデルを構造体で定義
type User struct { // typeで型の定義, structは構造体
	UserUuid    string `xorm:"varchar(36) pk" json:"userUUID"`                 // ユーザの一意のID
	Name        string `xorm:"varchar(25) not null" json:"name"`               // 名前
	UserType    string `xorm:"varchar(15) not null" json:"userType"`           // ユーザータイプ  // teacher, pupil, parents
	MailAddress string `xorm:"varchar(64) not null unique" json:"mailAddress"` // メアド
	Password    string `xorm:"varchar(60) not null" json:"password"`           // bcrypt化されたパスワード
	JtiUuid     string `xorm:"varchar(36) unique" json:"jwtUUID"`              // jwtクレームのuuid
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
func SaveJti(userId string, jti string) error {
	var user User // 取得したデータをマッピングする構造体
	// 更新前のレコードを取得
	if _, err := db.Where("user_uuid = ?", userId).Get(&user); err != nil {
		return err
	}

	// 受け取った新しい値を設定
	user.JtiUuid = jti

	// 更新を実行
	_, err := db.Update(&user) // 更新したレコードでテーブルの該当レコードを更新
	if err != nil {
		return err
	}
	return nil
}
