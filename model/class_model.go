package model

import (
	"time"
)

// クラステーブル
type Class struct {
	ClassUuid  string    `xorm:"varchar(36) pk" json:"classUUID"`       // クラスID
	ClassName  string    `xorm:"varchar(15) not null" json:"className"` // クラス名  // teacher, pupil, parents
	InviteCode string    `xorm:"char(4) unique" json:"inviteCode"`      // 招待ID
	ValidUntil time.Time `xorm:"datetime" json:"validUntil" `           // 有効期限
}

// テーブル名
func (Class) TableName() string {
	return "classes"
}

// テストデータ
func CreateClassTestData() {
	parsedTime := time.Now().AddDate(0, 0, -1)
	class1 := &Class{
		ClassUuid:  "09eba495-fe09-4f54-a856-9bea9536b661",
		ClassName:  "3-2 ふたば学級",
		InviteCode: "0000",
		ValidUntil: parsedTime,
	}
	db.Insert(class1)
	class2 := &Class{
		ClassUuid:  "817f600e-3109-47d7-ad8c-18b9d7dbdf8b",
		ClassName:  "つよつよガンギマリ塾 1-A",
		InviteCode: "9999",
		ValidUntil: parsedTime,
	}
	db.Insert(class2)
	class3 := &Class{
		ClassUuid:  "c085a423-d44b-4bde-a2ff-7e24736d33b5",
		ClassName:  "1-1 かかぽ学級",
		InviteCode: "1111",
		ValidUntil: parsedTime,
	}
	db.Insert(class3)

}

// クラス取得
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

// 複数クラス取得
func GetClasses(classUUids []string) ([]Class, error) {
	// 結果格納用変数
	var classes []Class

	// 全件取得
	err := db.In("class_uuid", classUUids).Find(&classes)
	if err != nil {
		return nil, err
	}

	return classes, nil
}

// 招待コードからクラスを取得
func GetClassByInviteCode(inviteCode string) (Class, error) {
	//結果格納用変数
	var class Class

	//inviteCodeで絞り込んで1件取得
	_, err := db.Where("invite_code =?", inviteCode).Get(
		&class,
	)

	// データが取得できなかったらerrを返す
	if err != nil {
		return Class{}, err
	}

	return class, nil

}

// 新規ユーザ登録
// 新しい構造体をレコードとして受け取り、usersテーブルにinsertし、成功した列数とerrorを返す
func CreateClass(record Class) (int64, error) {
	affected, err := db.Nullable("invite_code", "valid_until").Insert(record)
	return affected, err
}

// 招待コード更新
func UpdateInviteCode(record Class) (int64, error) {
	affected, err := db.Where("class_uuid = ?", record.ClassUuid).Cols("invite_code", "valid_until").Update(&record)
	return affected, err
}

func GetClassesByUUIDs(classUuids []string) ([]Class, error) {
	var classes []Class
	err := db.In("class_uuid", classUuids).Find(&classes)
	if err != nil {
		return nil, err
	}
	return classes, nil
}

// 期限の切れた招待コードと有効期限をnullにする
func DeleteExpiredInviteCodes() {
	_, err := db.Where("valid_until < ?", time.Now()).Nullable("invite_code", "valid_until").Update(&Class{})
	if err != nil {
		panic(err)
	}
}
