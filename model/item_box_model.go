package model

// 所持アイテムテーブル
type ItemBox struct {
	UserUuid string `xorm:"varchar(36) pk" json:"userUUID"`    // ユーザのUUID
	ItemUuid string `xorm:"varchar(36) pk" json:"nyariotUUID"` // アイテムUUID
	Quentity string `xorm:"int" json:"quentity"`               // アイテム所持数
}

// テーブル名
func (ItemBox) TableName() string {
	return "item_boxes"
}

// FK制約の追加
func InitItemBoxFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE item_boxes ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// ItemUuid
	_, err = db.Exec("ALTER TABLE item_boxes ADD FOREIGN KEY (item_uuid) REFERENCES items(item_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}
