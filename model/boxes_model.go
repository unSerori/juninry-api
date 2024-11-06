package model

type Box struct {
	HardwareUuid string `xorm:"varchar(36) pk" json:"hardwareUUID"`
	DepositPoint int    `xorm:"int default 0 not null" json:"depositPoint"`
}

// テーブル名
func (Box) TableName() string {
	return "boxes"
}

// FK制約の追加
func InitBoxesFK() error {
	_, err := db.Exec("ALTER TABLE boxes ADD FOREIGN KEY (hardware_uuid) REFERENCES hardwares(hardware_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateBoxesTestData() {
	boxes1 := &Box{
		HardwareUuid: "df2b1f4c-b49a-4068-80c5-3120dceb14c8",
	}
	db.Insert(boxes1)
}

// ボックスの現在のポイントを取得
func GetBoxDepositPoint(hardwareUUID string) (int, error) {
	// 結果格納用変数
	var box Box

	// 現在のポイントを取得
	_, err := db.Where("hardware_uuid = ?", hardwareUUID).Get(&box)
	if err != nil {
		return 0, err
	}
	return box.DepositPoint, nil
}

// ボックスのポイントを更新
func UpdateBoxDepositPoint(hardwareUUID string, depositPoint int) error {
	_, err := db.Where("hardware_uuid = ?", hardwareUUID).Update(&Box{DepositPoint: depositPoint})
	if err != nil {
		return err
	}
	return nil
}
