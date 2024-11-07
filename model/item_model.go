package model

// アイテムテーブル
type Item struct {
	ItemUuid      string `xorm:"varchar(36) pk" json:"itemUUID"`       // アイテムUUID
	ItemName      string `xorm:"varchar(25) " json:"itemName"`         // アイテム名
	ImagePath     string `xorm:"varchar(256) unique" json:"imagePath"` // アイテム画像パス
	Detail        string `xorm:"varchar(256)" json:"detail"`           // アイテム詳細
	Talk          string `xorm:"varchar(256)" json:"talk"`             // アイテム固有の会話
	SatityDegrees string `xorm:"int" json:"satityDegrees"`             // 空腹増加値
	Rarity        string `xorm:"int" json:"rarity"`                    // アイテムレアリティ
}

// テーブル名
func (Item) TableName() string {
	return "items"
}
