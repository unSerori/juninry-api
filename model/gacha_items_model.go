package model

// アイテムテーブル
type GachaItem struct {
	GachaNumber int    `xorm:"int" json:"gachaNumber"`         // ガチャ番号
	ItemUuid     string `xorm:"varchar(36) pk" json:"itemUUID"` // アイテムUUID
}

// テーブル名
func (GachaItem) TableName() string {
	return "gacha_items"
}

// アイテムデータ
func CreateGachaItemTestData() {

	items := []*GachaItem{
		{
			GachaNumber: 1,
			ItemUuid:     "563c5110-3441-4cb0-9764-f32c4385e975",
		},
		{
			GachaNumber: 2,
			ItemUuid:     "bf24b47d-a71b-4175-b5f1-fc622e5c5ac5",
		},
		{
			GachaNumber: 3,
			ItemUuid:     "0bb2b792-5d7c-449d-bb18-d96f8230508a",
		},
		{
			GachaNumber: 4,
			ItemUuid:     "5cdc139d-1f63-4695-a6e5-9d442f5d3fae",
		},
		{
			GachaNumber: 5,
			ItemUuid:     "5836d2df-3243-460b-b8f7-859a33516e1c",
		},
		{
			GachaNumber: 6,
			ItemUuid:     "4b52cf51-583e-4390-9f5c-5b7dbdfd65ef",
		},
		{
			GachaNumber: 7,
			ItemUuid:     "04de16b4-a7d2-4488-b609-0e5e3108bbc0",
		},
		{
			GachaNumber: 8,
			ItemUuid:     "05e50669-a2a7-428e-a243-95f3a2e5f98d",
		},
		{
			GachaNumber: 9,
			ItemUuid:     "cce677c2-423d-4049-9775-8dacc0d97cf1",
		},
		{
			GachaNumber: 10,
			ItemUuid:     "af9a9b0c-fbaf-4e4d-b0bb-0f2b9288ce50",
		},
		{
			GachaNumber: 11,
			ItemUuid:     "04de16b4-a7d2-4488-b609-0e5e3108bbc0",
		},
		{
			GachaNumber: 12,
			ItemUuid:     "f1e96b6d-1de4-4c82-8df5-2a10a1ab2f3f",
		},
		{
			GachaNumber: 13,
			ItemUuid:     "bb58c5e7-3a17-4987-b8e2-fb60b6d8979c",
		},
		{
			GachaNumber: 14,
			ItemUuid:     "5c2906b3-e759-4c5b-bd5d-3d848b91298e",
		},
		{
			GachaNumber: 15,
			ItemUuid:     "2d7b6377-3474-4a8c-a10c-b714ed3d9e11",
		},
	}

	// データベースに一括インサート
	for _, item := range items {
		db.Insert(item)
	}
}


func GetGachaItem() (GachaItem, error) {
	var gachaItem GachaItem
	_, err := db.Where("gacha_number BETWEEN ? AND ?", 1, 15).OrderBy("RAND()").Limit(1).Get(&gachaItem)
	if err!= nil {
		return gachaItem, err
	}

	return gachaItem, nil
}