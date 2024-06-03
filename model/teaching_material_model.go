package model

// 教材テーブル
type TeachingMaterial struct { // typeで型の定義, structは構造体
	TeachingMaterialUuid      string `xorm:"varchar(36) pk" json:"teachingMaterialUUID"`          // 教材ID
	TeachingMaterialName      string `xorm:"varchar(15) not null" json:"teachingMaterialName"`    // 教材名
	SubjectId                 int    `xorm:"not null" json:"subjectID"`                           // 教科ID
	TeachingMaterialImageUuid string `xorm:"varchar(36) unique" json:"teachingMaterialImageUUID"` // 教材画像
	ClassUuid                 string `xorm:"varchar(36) not null" json:"classUUID"`               // クラスID
}

// テーブル名
func (TeachingMaterial) TableName() string {
	return "teaching_materials"
}

// FK制約の追加
func InitTeachingMaterialFK() error {
	// SubjectId
	_, err := db.Exec("ALTER TABLE teaching_materials ADD FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// ClassUuid
	_, err = db.Exec("ALTER TABLE teaching_materials ADD FOREIGN KEY (class_uuid) REFERENCES classes(class_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateTeachingMaterialTestData() {
	tm1 := &TeachingMaterial{
		TeachingMaterialUuid:      "978f9835-5a16-4ac0-8581-7af8fac06b4e",
		TeachingMaterialName:      "漢字ドリル3",
		SubjectId:                 1,
		TeachingMaterialImageUuid: "a575f18c-d639-4b6d-ad57-a9d7a7f84575",
		ClassUuid:                 "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(tm1)
	tm2 := &TeachingMaterial{
		TeachingMaterialUuid:      "99cbb1be-5581-4607-b0ac-ab599edfd5d0",
		TeachingMaterialName:      "リピート1",
		SubjectId:                 4,
		TeachingMaterialImageUuid: "27fc9419-1673-4075-a73e-63ffa6c5d9f5",
		ClassUuid:                 "817f600e-3109-47d7-ad8c-18b9d7dbdf8b",
	}
	db.Insert(tm2)
}
