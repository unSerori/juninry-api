package model

// ユーザーの種類のテーブル
type UserType struct {
	UserTypeId int    `xorm:"pk autoincr" json:"userTypeId"`               // ユーザータイプID
	UserType   string `xorm:"varchar(15) not null unique" json:"userType"` // ユーザータイプ  // teacher, pupil, patron
}

// テーブル名
func (UserType) TableName() string {
	return "user_types"
}

// テストデータ
func CreateUserTypeTestData() {
	ut1 := &UserType{
		UserType: "teacher",
	}
	db.Insert(ut1)
	ut2 := &UserType{
		UserType: "pupil",
	}
	db.Insert(ut2)
	ut3 := &UserType{
		UserType: "parent",
	}
	db.Insert(ut3)
}
