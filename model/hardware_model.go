package model

// ハードウェアテーブル

type Hardware struct {
	HardwareUuid   string `xorm:"varchar(36) pk" json:"hardwareUUID"`
	HardwareTypeId int    `xorm:"int not null" json:"hardwareTypeId"`
}

// テーブル名
func (Hardware) TableName() string {
	return "hardwares"
}

// FK制約の追加
func InitHardwareFK() error {
	_, err := db.Exec("ALTER TABLE hardwares ADD FOREIGN KEY (hardware_type_id) REFERENCES hardware_types(hardware_type_id) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テスト用
func CreateHardwareTestData() {
	hardware1 := &Hardware{
		HardwareUuid:   "df2b1f4c-b49a-4068-80c5-3120dceb14c8",
		HardwareTypeId: 1,
	}
	db.Insert(hardware1)

	hardware2 := &Hardware{
		HardwareUuid:   "d611d471-5eb2-46a2-abaf-f758205f0d5f",
		HardwareTypeId: 1,
	}
	db.Insert(hardware2)
}
