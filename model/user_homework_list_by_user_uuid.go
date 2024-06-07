package model

import (
	"time"
)

// 実質Viewみたいな構造体
type UserHomework struct {
	HomeworkLimit             time.Time `xorm:"homework_limit"`               // 提出期限
	HomeworkUuid              string    `xorm:"homework_uuid"`                // 課題ID
	StartPage                 int       `xorm:"start_page"`                   // 開始ページ
	PageCount                 int       `xorm:"page_count"`                   // ページ数
	HomeworkNote              string    `xorm:"homework_note"`                // 課題の説明
	TeachingMaterialName      string    `xorm:"teaching_material_name"`       // 教材名
	SubjectId                 int       `json:"subjectID"`                    // 教科ID
	SubjectName               string    `json:"subjectName"`                  // 教科名
	TeachingMaterialImageUuid string    `xorm:"teaching_material_image_uuid"` // 画像ID どういう扱いになるのかな
	ClassName                 string    `xorm:"class_name"`                   // クラス名
}

func FindUserHomework(userUuid string) ([]UserHomework, error) {

	//クソデカ構造体のスライスを定義
	var userHomeworkList []UserHomework

	//クソデカ構造体をとるすごいやつだよ
	err := db.Table("homeworks").
		Join("LEFT", "teaching_materials", "homeworks.teaching_material_uuid = teaching_materials.teaching_material_uuid").
		// Join("LEFT", "users", "homeworks.homework_poster_uuid = users.user_uuid").	ユーザー名いらない子
		Join("LEFT", "subjects", "teaching_materials.subject_id = subjects.subject_id").
		Join("LEFT", "class_memberships", "teaching_materials.class_uuid = class_memberships.class_uuid").
		Join("LEFT", "classes", "teaching_materials.class_uuid = classes.class_uuid").
		Where("class_memberships.user_uuid = ?", userUuid).
		Select("homework_limit, homework_uuid, start_page, page_count, homework_note, teaching_material_name, subjects.subject_id, subject_name, teaching_material_image_uuid, class_name").
		OrderBy("teaching_materials.class_uuid").
		Find(&userHomeworkList)
	if err != nil { //エラーハンドル
		return nil, err
	}

	//クソデカ構造体のスライスを返す
	return userHomeworkList, nil
}
