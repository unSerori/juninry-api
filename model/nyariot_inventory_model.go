package model

// 所持ニャリオットテーブル
type NyariotInventory struct {
	UserUuid    string `xorm:"varchar(36) pk" json:"userUUID"`    // ユーザのUUID
	NyariotUuid string `xorm:"varchar(36) pk" json:"nyariotUUID"` // ニャリオットUUID
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
