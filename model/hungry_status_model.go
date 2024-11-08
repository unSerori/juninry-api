package model

import (
	"time"
)

// 満腹度テーブル
type HungryStatus struct {
	UserUuid      string    `xorm:"varchar(36) pk" json:"userUUID"`                                                          // ユーザのUUID
	SatityDegrees int       `xorm:"int not null default(100)" json:"satityDegrees"`                                          // 現在の空腹度
	NyariotUuid   string   `xorm:"varchar(36) not null default('c0768960-eb5f-4a60-8327-4171fd4b8a46')" json:"nyariotUUID"` // ニャリオットUUID
	LastGohanTime time.Time `xorm:"DATETIME not null" json:"lastGohanTime"`                                                      // 最後にご飯を食べた時間
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

func CreateHungryStatus(record HungryStatus) (int64, error) {
	affected, err := db.Insert(record)
	return affected, err
}

func ChangeNyariot(userUuid string, nyariotUuid string) (int64, error) {
		// 更新する構造体のインスタンスを作成
		changNyariot := HungryStatus{NyariotUuid: nyariotUuid}

	// ニャリオットを更新する
	affected, err := db.Where("user_uuid = ?", userUuid).Cols("nyariot_uuid").Update(&changNyariot)
	return affected, err
}

// ニャリオットの取得＋空腹度の更新
func GetMainNyairot(userUuid string, satityDegrees int) (HungryStatus, error) {
	var hungryStatuses HungryStatus
	_, err := db.Where("user_uuid = ?", userUuid).Get(&hungryStatuses)
	return hungryStatuses, err
}