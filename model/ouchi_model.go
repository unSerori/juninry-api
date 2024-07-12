package model

import (
	"time"
)

// おうちテーブル
type Ouchi struct {
	OuchiUuid  string    `xorm:"varchar(36) pk" json:"ouchiUUID"`       // ユーザータイプID
	OuchiName  string    `xorm:"varchar(15) not null" json:"ouchiName"` // ユーザータイプ  // teacher, pupil, parents
	InviteCode string    `xorm:"char(6) unique" json:"inviteCode"`      // 招待ID
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

// 新規おうち登録
// 新しい構造体をレコードとして受け取り、ouchiテーブルにinsertし、成功した列数とerrorを返す
func CreateOuchi(record Ouchi) (int64, error) {
	affected, err := db.Nullable("invite_code", "valid_until").Insert(record)
	return affected, err
}

// 招待コード更新
func UpdateOuchiInviteCode(record Ouchi) (int64, error) {
	affected, err := db.Where("ouchi_uuid = ?", record.OuchiUuid).Cols("invite_code", "valid_until").Update(&record)
	return affected, err
}

// おうち取得
func GetOuchi(ouchiUuid string) (Ouchi, error) {
	//結果格納用変数
	var ouchi Ouchi

	//ouchiUuidで絞り込んで1件取得
	_, err := db.Where("ouchi_uuid =?", ouchiUuid).Get(
		&ouchi,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return Ouchi{}, err
	}

	return ouchi, nil
}

// おうち招待コード取得
func GetOuchiInviteCode(inviteCode string) (Ouchi, error) {
	// 結果格納用
	var ouchi Ouchi

	//inviteCodeで絞り込んで1件取得
	_, err := db.Where("invite_code =?", inviteCode).Get(
		&ouchi,
	)

	// データが取得できなかったらerrを返す
	if err != nil {
		return Ouchi{}, err
	}

	return ouchi, nil
}

// 期限の切れたおうち招待コードと有効期限をnullにする
func DeleteExpiredOuchiInviteCodes() {
	_, err := db.Where("valid_until < ?", time.Now()).Nullable("invite_code", "valid_until").Update(&Ouchi{})
	if err != nil {
		panic(err)
	}
}

