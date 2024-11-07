package model

// アイテム画像管理テーブル
type ItemImage struct {
	ItemUuid  string `xorm:"varchar(36) pk" json:"itemUUID"`       //　アイテムID
	ImagePath string `xorm:"varchar(256) unique" json:"imagePath"` // アイテム画像パス
}

// テーブル名
func (ItemImage) TableName() string {
	return "item_images"
}

// テストデータ
func CreateItemImageTestData() {
	itemImage1 := &ItemImage{
		ItemUuid:  "563c5110-3441-4cb0-9764-f32c4385e97",
		ImagePath: "asset/images/item/IMG_0071.PNG",
	}
	db.Insert(itemImage1)
}

func GetImage(itemUuid string) {
	
}
