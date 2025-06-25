package model

// 所持アイテムテーブル
type ItemBox struct {
	UserUuid string `xorm:"varchar(36) pk" json:"userUUID"`    // ユーザのUUID
	ItemUuid string `xorm:"varchar(36) pk" json:"itemUUID"` // アイテムUUID
	Quantity int `xorm:"int" json:"quantity"`               // アイテム所持数
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

func GetUserItemBox(userUuid string, itemUuid string) (int, bool, error) {
	//結果格納用変数
	var itemBox ItemBox

	// userUuid, itemUuid で絞り込んだ結果
	found, err := db.Where("user_uuid = ? AND item_uuid = ?", userUuid, itemUuid).Get(&itemBox)

	// クエリ実行でエラーが発生した場合
	if err != nil {
		return 0, false, err
	}

	// アイテムが見つかった場合
	if found {
		return itemBox.Quantity, true, nil
	}

	// アイテムが見つからなかった場合
	return 0, false, nil
}

// アイテム減らす
func ReduceItemQuantity(userUuid string, itemUuid string, quantity int) (int64, error) {
	// 更新する構造体のインスタンスを作成
	itemQuantity := ItemBox{Quantity: quantity}

	// ニャリオットを更新する
	affected, err := db.Where("user_uuid = ? AND item_uuid = ?", userUuid, itemUuid).Cols("quantity").Update(&itemQuantity)
	return affected, err
}

//増やす
func UpdateItemQuantity(userUuid string, itmeUuid string) (int64, error) {

	affected, err := db.Where("user_uuid = ? AND item_uuid = ?", userUuid, itmeUuid).
		Incr("quantity", 1).
		Update(&ItemBox{})
	if err != nil {
		return 0, err
	}
	return affected, err
}

func CreateItemBox(record ItemBox) (int64, error) {
	affected, err := db.Insert(record)
	return affected, err
}