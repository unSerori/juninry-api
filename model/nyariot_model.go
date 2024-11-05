package model

// ニャリオットテーブル
type Nyariot struct {
	NyariotUuid      string `xorm:"varchar(36) pk" json:"nyariotUUID"`           // ニャリオットUUID
	NyariotName      string `xorm:"varchar(36) " json:"nyariotName"`             // ニャリオット名
	NyariotImagePath string `xorm:"varchar(256) unique" json:"nyariotImagePath"` // 画像パス
	Detail           string `xorm:"varchar(256)" json:"detail"`                  // 詳細
	Talk             string `xorm:"varchar(256)" json:"talk"`                    // ニャリオット固有の会話
}

// テーブル名
func (Nyariot) TableName() string {
	return "nyariots"
}
