package model

// アイテムテーブル
type GachaNyariot struct {
	GachaNumber int    `xorm:"int" json:"gachaNumber"`         // ガチャ番号
	NyariotUuid     string `xorm:"varchar(36) pk" json:"nyariotUUID"` // アイテムUUID
}

// テーブル名
func (GachaNyariot) TableName() string {
	return "gacha_nyariots"
}

// アイテムデータ
func CreateGachaNyariotTestData() {

	nyariots := []*GachaNyariot{
		{
			GachaNumber: 1,
			NyariotUuid:     "c0768960-eb5f-4a60-8327-4171fd4b8a46",
		},
		{
			GachaNumber: 2,
			NyariotUuid:     "baf8e173-0747-49d0-97c8-29a78e9319a9",
		},
		{
			GachaNumber: 3,
			NyariotUuid:     "ae30f602-9967-4851-b1e1-2ab10b1470bb",
		},
		{
			GachaNumber: 4,
			NyariotUuid:     "7b98eebc-7153-4903-9930-1b297bc5f120",
		},
		{
			GachaNumber: 5,
			NyariotUuid:     "9dbec0e8-8c9d-4901-a5d0-da952cbea1a4",
		},
		{
			GachaNumber: 6,
			NyariotUuid:     "5d42b7a8-348f-44c5-a364-17d77fcb9738",
		},
		{
			GachaNumber: 7,
			NyariotUuid:     "bf8e0c77-3397-48c5-8126-ebbb524e9ae3",
		},
		{
			GachaNumber: 8,
			NyariotUuid:     "f9539d0d-2853-4d6c-b3c3-9060af16eee3",
		},
		{
			GachaNumber: 9,
			NyariotUuid:     "9dbec0e8-8c9d-4901-a5d0-da952cbea1a4",
		},
		{
			GachaNumber: 10,
			NyariotUuid:     "cd47a101-8a5a-43af-bbd1-41c1d55a586e",
		},
	}

	// データベースに一括インサート
	for _, nyariot := range nyariots {
		db.Insert(nyariot)
	}
}

func GetGachaNyariot() (GachaNyariot, error) {
	var gachaNyariot GachaNyariot
	_, err := db.OrderBy("RAND()").Limit(1).Get(&gachaNyariot)
	if err!= nil {
		return gachaNyariot, err
	}

	return gachaNyariot, nil
}