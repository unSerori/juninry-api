package model

// ユーザのクラス所属中間テーブル
type ClassMembership struct {
	ClassUuid string `xorm:"varchar(36) pk" json:"classUUID"` // クラスID
	UserUuid  string `xorm:"varchar(36) pk" json:"userUUID"`  // ユーザーID
}

// テーブル名
func (ClassMembership) TableName() string {
	return "class_memberships"
}

// FK制約の追加
func InitClassMembershipFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE class_memberships ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	// ClassUuid
	_, err = db.Exec("ALTER TABLE class_memberships ADD FOREIGN KEY (class_uuid) REFERENCES classes(class_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateClassMembershipsTestData() {
	classMemberships1 := &ClassMembership{
		UserUuid:  "3cac1684-c1e0-47ae-92fd-6d7959759224",
		ClassUuid: "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(classMemberships1)
	classMemberships2 := &ClassMembership{
		UserUuid:  "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		ClassUuid: "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(classMemberships2)
}

// 用済みメソッド
// // user_uuidで絞り込み、所属構造体のスライスとerrorを返す
// func FindClassMemberships(userUuid string) ([]ClassMembership, error) {
// 	//ClassMemberships構造体のスライスを返すので定義
// 	var classMemberships []ClassMembership

// 	//uuidをWhere句で条件指定
// 	err := db.Where("user_uuid = ?", userUuid).Find(&classMemberships)
// 	if err != nil { //エラーハンドル
// 		return nil, err
// 	}

// 	//できたら返す
// 	return classMemberships, nil
// }
