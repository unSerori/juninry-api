package model

import (
	"time"
)

// 実質Viewみたいな構造体
type UserHomework struct {
	HomeworkLimit             time.Time // 提出期限
	HomeworkUuid              string    // 課題ID
	StartPage                 int       // 開始ページ
	PageCount                 int       // ページ数
	HomeworkNote              string    // 課題の説明
	TeachingMaterialName      string    // 教材名
	SubjectId                 int       // 教科ID
	SubjectName               string    // 教科名
	TeachingMaterialImageUuid string    // 画像ID どういう扱いになるのかな
	ClassName                 string    // クラス名
	SubmitFlag                int
}

// userUuidから課題データを取得、取得できなければエラーを返す
func FindUserHomework(userUuid string) ([]UserHomework, error) {
	//クソデカ構造体のスライスを定義
	var userHomeworkList []UserHomework

	//クソデカ構造体をとるすごいやつだよ
	err := db.Table("homeworks").
		Join("LEFT", "teaching_materials", "homeworks.teaching_material_uuid = teaching_materials.teaching_material_uuid").
		Where("homework_limit > ?", time.Now()).
		Join("LEFT", "subjects", "teaching_materials.subject_id = subjects.subject_id").
		Join("LEFT", "class_memberships", "teaching_materials.class_uuid = class_memberships.class_uuid").
		Join("LEFT", "classes", "teaching_materials.class_uuid = classes.class_uuid").
		Join("LEFT", "homework_submissions", "homeworks.homework_uuid = homework_submissions.homework_uuid AND homework_submissions.user_uuid = ?", userUuid).
		Where("class_memberships.user_uuid = ?", userUuid).
		Select("homework_limit, homeworks.homework_uuid, start_page, page_count, homework_note, teaching_material_name, subjects.subject_id, subject_name, teaching_material_image_uuid, class_name, if(homework_submissions.user_uuid IS NOT NULL, 1, 0) as submit_flag").
		OrderBy("homework_limit, teaching_materials.class_uuid, submit_flag").
		Find(&userHomeworkList)
	if err != nil { //エラーハンドル ただエラー投げてるだけ
		return nil, err
	}

	//クソデカ構造体のスライスを返す
	return userHomeworkList, nil
}

// userUuidから期限を次の日に指定して課題を取得、取得できなければエラーを返す
func FindUserHomeworkforNextday(userUuid []string) ([]UserHomework, error) {
	//クソデカ構造体のスライスを定義
	var userHomeworkList []UserHomework
	// 現在の日付を取得
	now := time.Now()
	// 現在の日付から時間部分を0時にリセットし、次の日を取得
	tomorrowMidnight := time.Date(
		now.Year(), now.Month(), now.Day()+2, 0, 0, 0, 0, now.Location(),
	)

	// userUuid をインターフェース型にキャスト
	uuidInterfaces := make([]interface{}, len(userUuid))
	for i, v := range userUuid {
		uuidInterfaces[i] = v
	}

	//クソデカ構造体をとるすごいやつだよ
	err := db.Table("homeworks").
		Join("LEFT", "teaching_materials", "homeworks.teaching_material_uuid = teaching_materials.teaching_material_uuid").
		Where("homework_limit BETWEEN ? AND ?", now,tomorrowMidnight).	// 今日の日付以降
		Join("LEFT", "subjects", "teaching_materials.subject_id = subjects.subject_id").
		Join("LEFT", "class_memberships", "teaching_materials.class_uuid = class_memberships.class_uuid").
		Join("LEFT", "classes", "teaching_materials.class_uuid = classes.class_uuid").
		Join("LEFT", "homework_submissions", "homeworks.homework_uuid = homework_submissions.homework_uuid AND homework_submissions.user_uuid IN (?)", uuidInterfaces...).
		In("class_memberships.user_uuid", uuidInterfaces...).
		Select("homework_limit, homeworks.homework_uuid, start_page, page_count, homework_note, teaching_material_name, subjects.subject_id, subject_name, teaching_material_image_uuid, class_name, if(homework_submissions.user_uuid IS NOT NULL, 1, 0) as submit_flag").
		OrderBy("homework_limit, teaching_materials.class_uuid, submit_flag").
		Find(&userHomeworkList)
	if err != nil { //エラーハンドル ただエラー投げてるだけ
		return nil, err
	}

	//クソデカ構造体のスライスを返す
	return userHomeworkList, nil
}
