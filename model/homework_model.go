package model

import (
	"fmt"
	"time"
)

// 課題テーブル
type Homework struct { // typeで型の定義, structは構造体
	HomeworkUuid         string    `xorm:"varchar(36) pk" json:"taskUUID"`                   // タスクのID
	HomeworkLimit        time.Time `xorm:"DATETIME not null" json:"taskLimit"`               // タスクの期限
	TeachingMaterialUuid string    `xorm:"varchar(36) not null" json:"teachingMaterialUUID"` // 教材ID
	StartPage            int       `json:"startPage"`                                        // 開始ページ
	PageCount            int       `xorm:"not null" json:"pageCount"`                        // ページ数
	HomeworkPosterUuid   string    `xorm:"varchar(36) not null" json:"homeworkPosterUUID"`   // 投稿者ID
	HomeworkNote         string    `xorm:"varchar(255)" json:"homeworkNote"`                 // 宿題説明
}

// テーブル名
func (Homework) TableName() string {
	return "homework"
}

// FK制約の追加
func InitHomeworkFK() error {
	// HomeworkPosterUuid
	_, err := db.Exec("ALTER TABLE homework ADD FOREIGN KEY (homework_poster_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateHomeworkTestData() {
	parsedTime, _ := time.Parse(time.RFC3339, "2024-06-03 06:14:11.515967422 +0000 UTC m=+0.318201036")
	fmt.Println("homework_limit: ", parsedTime)
	homework1 := &Homework{
		HomeworkUuid:         "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		HomeworkLimit:        parsedTime,
		TeachingMaterialUuid: "978f9835-5a16-4ac0-8581-7af8fac06b4e",
		StartPage:            24,
		PageCount:            2,
		HomeworkPosterUuid:   "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		HomeworkNote:         "がんばってくださ～い＾＾",
	}
	db.Insert(homework1)
}
