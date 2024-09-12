package model

import (
	"juninry-api/domain"
	"juninry-api/utility/custom"
)

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

// エンティティとの相互変換

// テーブルモデルをドメインエンティティに変換
func (tm *TeachingMaterial) ToDomainEntity() *domain.TeachingMaterial {
	return &domain.TeachingMaterial{
		TeachingMaterialUuid:      tm.TeachingMaterialUuid,
		TeachingMaterialName:      tm.TeachingMaterialName,
		SubjectId:                 tm.SubjectId,
		TeachingMaterialImageUuid: tm.TeachingMaterialImageUuid,
		ClassUuid:                 tm.ClassUuid,
	}
}

// ドメインエンティティをテーブルモデルに変換
func FromDomainEntity(de *domain.TeachingMaterial) *TeachingMaterial {
	return &TeachingMaterial{
		TeachingMaterialUuid:      de.TeachingMaterialUuid,
		TeachingMaterialName:      de.TeachingMaterialName,
		SubjectId:                 de.SubjectId,
		TeachingMaterialImageUuid: de.TeachingMaterialImageUuid,
		ClassUuid:                 de.ClassUuid,
	}
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
	tm3 := &TeachingMaterial{
		TeachingMaterialUuid:      "22b78a9d-cfc2-4f0e-bb2f-19002dd259f3",
		TeachingMaterialName:      "せいかつ",
		SubjectId:                 3,
		TeachingMaterialImageUuid: "4391c3e9-0151-45e8-ae70-d20879dacc95",
		ClassUuid:                 "817f600e-3109-47d7-ad8c-18b9d7dbdf8b",
	}
	db.Insert(tm3)
	tm4 := &TeachingMaterial{
		TeachingMaterialUuid:      "978f9835-5a16-4ac0-8581-7affac06b4e",
		TeachingMaterialName:      "計算ドリル",
		SubjectId:                 2,
		TeachingMaterialImageUuid: "a575f18c-d639-4b6d-ad5-a9d7a7f84575",
		ClassUuid:                 "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(tm4)
	tm5 := &TeachingMaterial{
		TeachingMaterialUuid:      "978f9835-5a16-4ac0-8581-7af8fac0b4e",
		TeachingMaterialName:      "理科ワーク",
		SubjectId:                 3,
		TeachingMaterialImageUuid: "a575f18c-d639-4b6d-ad57a9d7a7f84575",
		ClassUuid:                 "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(tm5)
}

// クラスIDから教材一覧を取得
func FindTeachingMaterials(classUuids []string) ([]TeachingMaterial, error) {
	var teachingMaterials []TeachingMaterial
	err := db.In("class_uuid", classUuids).Find(&teachingMaterials)
	if err != nil {
		return nil, err
	}
	return teachingMaterials, nil
}

// tmIdから教材がどのクラスで発行されたものか取得
func GetClassId(tmId string) (string, error) {
	var tm TeachingMaterial // 取得したデータをマッピングする構造体
	isFound, err := db.Where("teaching_material_uuid = ?", tmId).Get(&tm)
	if err != nil {
		return "", err
	}
	if !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`!isFound`で、影響を与えるSQL文の時は`affected == 0`でハンドリング
		return "", custom.NewErr(custom.ErrTypeNoFoundR)
	}

	return tm.ClassUuid, nil

}

// 教材名取得
func GetTmName(tmId string) (string, error) {
	var tm TeachingMaterial // 取得したデータをマッピングする構造体
	isFound, err := db.Where("teaching_material_uuid = ?", tmId).Get(&tm)
	if err != nil {
		return "", err
	}
	if !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`!isFound`で、影響を与えるSQL文の時は`affected == 0`でハンドリング
		return "", custom.NewErr(custom.ErrTypeNoFoundR)
	}

	return tm.TeachingMaterialName, nil
}

// 教科IDを取得
func GetSubjectId(tmId string) (int, error) {
	var tm TeachingMaterial // 取得したデータをマッピングする構造体
	isFound, err := db.Where("teaching_material_uuid = ?", tmId).Get(&tm)
	if err != nil {
		return 0, err
	}
	if !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`!isFound`で、影響を与えるSQL文の時は`affected == 0`でハンドリング
		return 0, custom.NewErr(custom.ErrTypeNoFoundR)
	}

	return tm.SubjectId, nil
}
