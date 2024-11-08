package model

// 所持ニャリオットテーブル
type NyariotInventory struct {
	UserUuid     string `xorm:"varchar(36) pk" json:"userUUID"`    // ユーザのUUID
	NyariotUuid  string `xorm:"varchar(36) pk" json:"nyariotUUID"` // ニャリオットUUID
	ConvexNumber int    `xorm:"int" json:"convexNumber"`           // 凸数
}

// テーブル名
func (NyariotInventory) TableName() string {
	return "nyariot_inventories"
}

// FK制約の追加
func InitNyariotInventoryFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE nyariot_inventories ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// NyariotUuid
	_, err = db.Exec("ALTER TABLE nyariot_inventories ADD FOREIGN KEY (nyariot_uuid) REFERENCES nyariots(nyariot_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

func CreateNyariotInventory(record NyariotInventory) (int64, error) {
	affected, err := db.Insert(record)
	return affected, err
}

func GetUserNyariotInbentory(userUuid string, nyariotUuid string) (int, bool, error) {
	//結果格納用変数
	var nyariotInventory NyariotInventory

	// userUuid, itemUuid で絞り込んだ結果
	found, err := db.Where("user_uuid = ? AND nyariot_uuid = ?", userUuid, nyariotUuid).Get(&nyariotInventory)

	// クエリ実行でエラーが発生した場合
	if err != nil {
		return 0, false, err
	}

	// アイテムが見つかった場合
	if found {
		return nyariotInventory.ConvexNumber, true, nil
	}

	// アイテムが見つからなかった場合
	return 0, false, nil
}
