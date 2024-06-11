package model

// クラステーブル
type Class struct {
	ClassUuid string `xorm:"varchar(36) pk" json:"classUUID"`       // ユーザータイプID
	ClassName string `xorm:"varchar(15) not null" json:"className"` // ユーザータイプ  // teacher, pupil, parents
}

// テーブル名
func (Class) TableName() string {
	return "classes"
}

// テストデータ
func CreateClassTestData() {
	class1 := &Class{
		ClassUuid: "09eba495-fe09-4f54-a856-9bea9536b661",
		ClassName: "3-2 ふたば学級",
	}
	db.Insert(class1)
	class2 := &Class{
		ClassUuid: "817f600e-3109-47d7-ad8c-18b9d7dbdf8b",
		ClassName: "つよつよガンギマリ塾 1-A",
	}
	db.Insert(class2)
}

//クラス取得
func GetClass(classUuid string) (Class, error) {
	//結果格納用変数
	var class Class

	//classUuidで絞り込んで1件取得
	_, err := db.Where("class_uuid =?", classUuid).Get(
		&class,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return Class{}, err
	}

	return class, nil
}
