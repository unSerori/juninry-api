package model

type Box struct {
	HardwareUuid string `xorm:"varchar(36) pk" json:"hardwareUUID"`
	DepositPoint int    `xorm:"int default 0 not null" json:"depositPoint"`
	BoxStatus    int    `xorm:"int default 0 not null" json:"boxStatus"` // 0: 何も登録されていない状態　1: ポイント貯めたりできる状態　2: メンテナンス中
	OuchiUuid    string `xorm:"varchar(36) not null" json:"ouchiUUID"`
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

	_, err = db.Exec("ALTER TABLE boxes ADD FOREIGN KEY (ouchi_uuid) REFERENCES ouchies(ouchi_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateBoxesTestData() {
	boxes1 := &Box{
		HardwareUuid: "df2b1f4c-b49a-4068-80c5-3120dceb14c8",
		OuchiUuid:    "2e17a448-985b-421d-9b9f-62e5a4f28c49",
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

// ボックスを取得
func GetBox(hardwareUUID string) (Box, error) {
	// 結果格納用変数
	var box Box

	// 現在のポイントを取得
	_, err := db.Where("hardware_uuid = ?", hardwareUUID).Get(&box)
	if err != nil {
		return Box{}, err
	}
	return box, nil
}

// ボックスの一覧を取得
func GetBoxes(ouchiUuid string) ([]Box, error) {
	// 結果格納用変数
	var boxes []Box

	// 現在のポイントを取得
	err := db.Where("ouchi_uuid = ?", ouchiUuid).Find(&boxes)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

// ボックスのポイントを更新
func UpdateBoxDepositPoint(hardwareUUID string, depositPoint int) error {
	_, err := db.Where("hardware_uuid = ?", hardwareUUID).Update(&Box{DepositPoint: depositPoint})
	if err != nil {
		return err
	}
	return nil
}
