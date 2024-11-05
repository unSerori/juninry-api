package model

import(
	"time"
)

// スタンプテーブル
type Stamp struct {
	UserUuid string `xorm:"varchar(36) pk" json:"userUUID"` // ユーザのUUID
	Quentity string `xorm:"int" json:"quentity"`            // スタンプの数
	LastLoginTime time.Time `xorm:"DATETIME not null" json:"lastLoginTime"` // 最後にログインした時間
}

// テーブル名
func (Stamp) TableName() string {
	return "stamps"
}

// FK制約の追加
func InitStampFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE stamps ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}