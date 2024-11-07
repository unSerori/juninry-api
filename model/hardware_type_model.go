package model

// ハードウェアの種類のテーブル
type HardwareType struct {
	HardwareTypeId   int    `xorm:"not null pk autoincr" json:"hardwareTypeId"`
	HardwareTypeName string `xorm:"varchar(30) not null unique" json:"name"`
}

// テーブル名
func (HardwareType) TableName() string {
	return "hardware_types"
}



// テストデータ
func CreateHardwareTypeTestData() {
	hardwareType1 := &HardwareType{
		HardwareTypeName: "宝箱",
	}
	db.Insert(hardwareType1)
}
