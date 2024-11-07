package model

// ニャリオットテーブル
type Nyariot struct {
	NyariotUuid      string `xorm:"varchar(36) pk" json:"nyariotUUID"`           // ニャリオットUUID
	NyariotName      string `xorm:"varchar(36) " json:"nyariotName"`             // ニャリオット名
	NyariotImagePath string `xorm:"varchar(256) unique" json:"nyariotImagePath"` // 画像パス
	NyariotNumber    int    `xorm:"int" json:"nyariotNumber"`                    // ニャリオット番号
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
		NyariotName:      "1位ニャリオット",
		NyariotImagePath: "asset/images/nyariot/IMG_0066.PNG",
		NyariotNumber:    1,
		Detail:           "かけっこで1位になったみたい",
		Talk:             "やったーーーー",
	}
	db.Insert(nyariot1)
}
