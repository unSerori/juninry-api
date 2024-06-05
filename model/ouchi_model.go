package model

// おうちテーブル
type Ouchi struct {
	OuchiUuid string `xorm:"varchar(36) pk" json:"ouchiUUID"`       // ユーザータイプID
	OuchiName string `xorm:"varchar(15) not null" json:"ouchiName"` // ユーザータイプ  // teacher, pupil, parents
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
