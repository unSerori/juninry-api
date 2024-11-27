package model

// アイテムテーブル
type Item struct {
	ItemUuid      string `xorm:"varchar(36) pk" json:"itemUUID"`            // アイテムUUID
	ItemName      string `xorm:"varchar(25) " json:"itemName"`              // アイテム名
	ImagePath     string `xorm:"varchar(256) unique" json:"imagePath"`      // アイテム画像パス
	ItemNumber    int    `xorm:"int not null default(0)" json:"itemNumber"` // アイテム番号
	Detail        string `xorm:"varchar(256)" json:"detail"`                // アイテム詳細
	Talk          string `xorm:"varchar(256)" json:"talk"`                  // アイテム固有の会話
	SatityDegrees int    `xorm:"int" json:"satityDegrees"`                  // 空腹増加値
	Rarity        int    `xorm:"int" json:"rarity"`                         // アイテムレアリティ 1:N 2:R 3:SR
}

// テーブル名
func (Item) TableName() string {
	return "items"
}

// アイテムデータ
func CreateItemTestData() {
	item1 := &Item{
		ItemUuid:      "563c5110-3441-4cb0-9764-f32c4385e975",
		ItemName:      "おさかな",
		ImagePath:     "asset/images/item/IMG_0071.PNG",
		ItemNumber:    1,
		Detail:        "ちからもスピードもほとんどダメ",
		Talk:          "遥か大昔はそれなりに強かったみたい",
		SatityDegrees: 10,
		Rarity:        1,
	}
	db.Insert(item1)

	item2 := &Item{
		ItemUuid:      "bf24b47d-a71b-4175-b5f1-fc622e5c5ac5",
		ItemName:      "おいしいみず",
		ImagePath:     "asset/images/item/water.PNG",
		ItemNumber:    2,
		Detail:        "ミネラルたっぷりのみず",
		Talk:          "しれんサポーターに渡すとビビリだまが貰えるよ",
		SatityDegrees: 20,
		Rarity:        2,
	}
	db.Insert(item2)

	item3 := &Item{
		ItemUuid:      "0bb2b792-5d7c-449d-bb18-d96f8230508a",
		ItemName:      "ワイヤー",
		ImagePath:     "asset/images/item/wiyer.PNG",
		ItemNumber:    3,
		Detail:        "お尻から出す糸はワイヤーに匹敵する強度",
		Talk:          "強さの秘密が研究されているよ",
		SatityDegrees: 50,
		Rarity:        3,
	}
	db.Insert(item3)

	item4 := &Item{
		ItemUuid:      "5cdc139d-1f63-4695-a6e5-9d442f5d3fae",
		ItemName:      "ホイップクリーム",
		ImagePath:     "asset/images/item/IMG_01.PNG",
		ItemNumber:    4,
		Detail:        "デコレーションはダイヤモンド並みの強度",
		Talk:          "深みのあるあまさで食べたみんなを幸せにするよ",
		SatityDegrees: 20,
		Rarity:        2,
	}
	db.Insert(item4)

	item5 := &Item{
		ItemUuid:      "5836d2df-3243-460b-b8f7-859a33516e1c",
		ItemName:      "クリーム",
		ImagePath:     "asset/images/item/IMG_02.PNG",
		ItemNumber:    5,
		Detail:        "デコレみの強度",
		Talk:          "たみんなを幸せにするよ",
		SatityDegrees: 20,
		Rarity:        2,
	}
	db.Insert(item5)

	item6 := &Item{
		ItemUuid:      "4b52cf51-583e-4390-9f5c-5b7dbdfd65ef",
		ItemName:      "あああああああ",
		ImagePath:     "asset/images/item/IMG_3.PNG",
		ItemNumber:    6,
		Detail:        "いいいいい",
		Talk:          "ううううううう",
		SatityDegrees: 10,
		Rarity:        1,
	}
	db.Insert(item6)

	item7 := &Item{
		ItemUuid:      "04de16b4-a7d2-4488-b609-0e5e3108bbc0",
		ItemName:      "132453",
		ImagePath:     "asset/images/item/IMG_7.PNG",
		ItemNumber:    7,
		Detail:        "hayakunetai",
		Talk:          "turatanienn",
		SatityDegrees: 10,
		Rarity:        1,
	}
	db.Insert(item7)

	item8 := &Item{
		ItemUuid:      "05e50669-a2a7-428e-a243-95f3a2e5f98d",
		ItemName:      "あああああああ",
		ImagePath:     "asset/images/item/IMG_8.PNG",
		ItemNumber:    8,
		Detail:        "hutonngakoisii",
		Talk:          "makuragakoisii",
		SatityDegrees: 10,
		Rarity:        1,
	}
	db.Insert(item8)

	item9 := &Item{
		ItemUuid:      "cce677c2-423d-4049-9775-8dacc0d97cf1",
		ItemName:      "67890",
		ImagePath:     "asset/images/item/IMG_9.PNG",
		ItemNumber:    9,
		Detail:        "iiiiii",
		Talk:          "uuuuuu",
		SatityDegrees: 20,
		Rarity:        2,
	}
	db.Insert(item9)

	item10 := &Item{
		ItemUuid:      "af9a9b0c-fbaf-4e4d-b0bb-0f2b9288ce50",
		ItemName:      "12345",
		ImagePath:     "asset/images/item/IMG_10.PNG",
		ItemNumber:    10,
		Detail:        "qwert",
		Talk:          "qwert",
		SatityDegrees: 10,
		Rarity:        1,
	}
	db.Insert(item10)

	items := []*Item{
		{
			ItemUuid:      "04de16b4-a7d2-4488-b609-0e5e3108bbc0",
			ItemName:      "Item 11",
			ImagePath:     "asset/images/item/IMG_11.PNG",
			ItemNumber:    11,
			Detail:        "Detail of item 11",
			Talk:          "Talk for item 11",
			SatityDegrees: 20,
			Rarity:        2,
		},
		{
			ItemUuid:      "f1e96b6d-1de4-4c82-8df5-2a10a1ab2f3f",
			ItemName:      "Item 12",
			ImagePath:     "asset/images/item/IMG_12.PNG",
			ItemNumber:    12,
			Detail:        "Detail of item 12",
			Talk:          "Talk for item 12",
			SatityDegrees: 10,
			Rarity:        1,
		},
		{
			ItemUuid:      "bb58c5e7-3a17-4987-b8e2-fb60b6d8979c",
			ItemName:      "Item 13",
			ImagePath:     "asset/images/item/IMG_13.PNG",
			ItemNumber:    13,
			Detail:        "Detail of item 13",
			Talk:          "Talk for item 13",
			SatityDegrees: 50,
			Rarity:        3,
		},
		{
			ItemUuid:      "5c2906b3-e759-4c5b-bd5d-3d848b91298e",
			ItemName:      "Item 14",
			ImagePath:     "asset/images/item/IMG_14.PNG",
			ItemNumber:    14,
			Detail:        "Detail of item 14",
			Talk:          "Talk for item 14",
			SatityDegrees: 10,
			Rarity:        1,
		},
		{
			ItemUuid:      "2d7b6377-3474-4a8c-a10c-b714ed3d9e11",
			ItemName:      "Item 15",
			ImagePath:     "asset/images/item/IMG_15.PNG",
			ItemNumber:    15,
			Detail:        "Detail of item 15",
			Talk:          "Talk for item 15",
			SatityDegrees: 20,
			Rarity:        2,
		},
	}

	// データベースに一括インサート
	for _, item := range items {
		db.Insert(item)
	}

}

// 全アイテム取得
func GetItems() ([]Item, error) {
	// 結果を格納する変数宣言(findの結果)
	var items []Item

	//データを全件取得
	err := db.Asc("item_number").Find(&items)

	// データが取得できなかったらerrを返す
	if err != nil {
		return nil, err
	}

	// エラーが出なければ取得結果を返す
	return items, nil
}

func GetItem(itemUuid string) (Item, error) {
	// 結果を格納する変数宣言
	var item Item

	_, err := db.Where("item_uuid = ?", itemUuid).Get(&item)
	if err != nil {
		return Item{}, err
	}

	return item, nil

}
