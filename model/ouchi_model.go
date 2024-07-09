package model

import (
	"time"
)

// おうちテーブル
type Ouchi struct {
	OuchiUuid string `xorm:"varchar(36) pk" json:"ouchiUUID"`       // ユーザータイプID
	OuchiName string `xorm:"varchar(15) not null" json:"ouchiName"` // ユーザータイプ  // teacher, pupil, parents
	InviteCode string    `xorm:"char(4) unique" json:"inviteCode"`      // 招待ID
	ValidUntil time.Time `xorm:"datetime" json:"validUntil" `           // 有効期限
}

// テーブル名
func (Ouchi) TableName() string {
	return "ouchies"
}

// テストデータ
func CreateOuchiTestData() {
	ouchi1 := &Ouchi{
		OuchiUuid: "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		OuchiName: "piyonaka家",
	}
	db.Insert(ouchi1)
	ouchi2 := &Ouchi{
		OuchiUuid: "48743657-b250-42f7-8850-64c430a980ba",
		OuchiName: "たんぽぽ施設",
	}
	db.Insert(ouchi2)
}

// 招待コード更新
func UpdateOuchiInviteCode(record Ouchi) (int64, error) {
	affected, err := db.Where("ouchi_uuid = ?", record.OuchiUuid).Cols("invite_code", "valid_until").Update(&record)
	return affected, err
}
