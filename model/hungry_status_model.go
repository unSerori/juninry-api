package model

import (
	"time"
)

// スタンプテーブル
type HungryStatus struct {
	UserUuid      string    `xorm:"varchar(36) pk" json:"userUUID"`         // ユーザのUUID
	SatityDegrees string    `xorm:"int" json:"satityDegrees"`               // 現在の空腹度
	NyariotUuid   string    `xorm:"varchar(36)" json:"nyariotUUID"`         // ニャリオットUUID
	LastGohanTime time.Time `xorm:"DATETIME not null" json:"lastGohanTime"` // 最後にご飯を食べた時間
}

// テーブル名
func (HungryStatus) TableName() string {
	return "hungry_statuses"
}

// FK制約の追加
func InitHungryStatusFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE hungry_statuses ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// NyariotUuid
	_, err = db.Exec("ALTER TABLE hungry_statuses ADD FOREIGN KEY (nyariot_uuid) REFERENCES nyariots(nyariot_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}
