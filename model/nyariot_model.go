package model

// ニャリオットテーブル
type Nyariot struct {
	NyariotUuid      string `xorm:"varchar(36) pk" json:"nyariotUUID"`           // ニャリオットUUID
	NyariotName      string `xorm:"varchar(36) " json:"nyariotName"`             // ニャリオット名
	NyariotImagePath string `xorm:"varchar(256) unique" json:"nyariotImagePath"` // 画像パス
	Nyarindex        int    `xorm:"int" json:"nyariotNumber"`                    // ニャリオット番号
	Detail           string `xorm:"varchar(256)" json:"detail"`                  // 詳細
	Talk             string `xorm:"varchar(256)" json:"talk"`                    // ニャリオット固有の会話
}

// テーブル名
func (Nyariot) TableName() string {
	return "nyariots"
}

// アイテムデータ
func CreateNyariotTestData() {
	nyariot1 := &Nyariot{
		NyariotUuid:      "c0768960-eb5f-4a60-8327-4171fd4b8a46",
		NyariotName:      "デフォルトニャリオット",
		NyariotImagePath: "nyariotImage/cat.PNG",
		Nyarindex:        1,
		Detail:           "野生のニャリオットが飛び出してきた",
		Talk:             "初めましてニャ",
	}
	db.Insert(nyariot1)

	nyariot2 := &Nyariot{
		NyariotUuid:      "baf8e173-0747-49d0-97c8-29a78e9319a9",
		NyariotName:      "じゃんけんニャリオット",
		NyariotImagePath: "nyariotImage/IMG_0067.PNG",
		Nyarindex:        2,
		Detail:           "じゃんけんニャリオットが勝負を仕掛けてきた！",
		Talk:             "君はどれを出す？",
	}
	db.Insert(nyariot2)

	nyariot3 := &Nyariot{
		NyariotUuid:      "ae30f602-9967-4851-b1e1-2ab10b1470bb",
		NyariotName:      "きだるげニャリオット",
		NyariotImagePath: "nyariotImage/IMG_0068.PNG",
		Nyarindex:        3,
		Detail:           "今日は雨の日",
		Talk:             "やる気が出ないよ～",
	}
	db.Insert(nyariot3)

	nyariot4 := &Nyariot{
		NyariotUuid:      "7b98eebc-7153-4903-9930-1b297bc5f120",
		NyariotName:      "1位ニャリオット",
		NyariotImagePath: "nyariotImage/IMG_0066.PNG",
		Nyarindex:        1,
		Detail:           "かけっこで1位になったみたい",
		Talk:             "やったーーーー",
	}
	db.Insert(nyariot4)

	nyariot5 := &Nyariot{
		NyariotUuid:      "9dbec0e8-8c9d-4901-a5d0-da952cbea1a4",
		NyariotName:      "Bニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0070.PNG",
		Nyarindex:        5,
		Detail:           "BBBBBBBBB",
		Talk:             "BBBBBBBB",
	}
	db.Insert(nyariot5)

	nyariot6 := &Nyariot{
		NyariotUuid:      "5d42b7a8-348f-44c5-a364-17d77fcb9738",
		NyariotName:      "Cニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0071.PNG",
		Nyarindex:        6,
		Detail:           "cccccccc",
		Talk:             "cccccccc",
	}
	db.Insert(nyariot6)

	nyariot7 := &Nyariot{
		NyariotUuid:      "bf8e0c77-3397-48c5-8126-ebbb524e9ae3",
		NyariotName:      "dニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0072.PNG",
		Nyarindex:        7,
		Detail:           "dddddddd",
		Talk:             "dddddddd",
	}
	db.Insert(nyariot7)

	nyariot8 := &Nyariot{
		NyariotUuid:      "f9539d0d-2853-4d6c-b3c3-9060af16eee3",
		NyariotName:      "Eニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0073.PNG",
		Nyarindex:        8,
		Detail:           "EEEEEEEE",
		Talk:             "EEEEEEEE",
	}
	db.Insert(nyariot8)

	nyariot9 := &Nyariot{
		NyariotUuid:      "9dbec0e8-8c9d-4901-a5d0-da952cbea1a4",
		NyariotName:      "Fニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0074.PNG",
		Nyarindex:        9,
		Detail:           "FFFFFFFFF",
		Talk:             "ffffffffff",
	}
	db.Insert(nyariot9)

	nyariot10 := &Nyariot{
		NyariotUuid:      "cd47a101-8a5a-43af-bbd1-41c1d55a586e",
		NyariotName:      "Gニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0075.PNG",
		Nyarindex:        10,
		Detail:           "gGGGGGGGGgg",
		Talk:             "GGGGGGGG",
	}
	db.Insert(nyariot10)
}

// 全ニャリオット取得
func GetNyariots() ([]Nyariot, error) {
	// 結果格納用変数
	var nyariot []Nyariot

	// データ全件取得
	err := db.Asc("nyarindex").Find(&nyariot)
	// データが取得できなかったらerrを返す
	if err != nil {
		return nil, err
	}

	// エラーが出なければ取得結果を返す
	return nyariot, nil
}

func GetNyariot(NyariotUuid string) (Nyariot, error) {
		// 結果を格納する変数宣言
		var nyariot Nyariot

		_, err := db.Where("nyariot_uuid = ?", NyariotUuid).Get(&nyariot)
		if err != nil {
			return Nyariot{}, err
		}
	
		return nyariot, nil
}
