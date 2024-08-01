package model

import (
	"errors"
	"fmt"
	"juninry-api/utility/custom"
	"time"

	"github.com/go-sql-driver/mysql"
)

// 課題テーブル
type Homework struct { // typeで型の定義, structは構造体
	HomeworkUuid         string    `xorm:"varchar(36) pk" json:"homeworkUUID"`               // タスクのID
	HomeworkLimit        time.Time `xorm:"DATETIME not null" json:"homeworkLimit"`           // タスクの期限
	TeachingMaterialUuid string    `xorm:"varchar(36) not null" json:"teachingMaterialUUID"` // 教材ID
	StartPage            int       `json:"startPage"`                                        // 開始ページ
	PageCount            int       `xorm:"not null" json:"pageCount"`                        // ページ数
	HomeworkPosterUuid   string    `xorm:"varchar(36) not null" json:"homeworkPosterUUID"`   // 投稿者ID
	HomeworkNote         string    `xorm:"varchar(255)" json:"homeworkNote"`                 // 宿題説明
}

// テーブル名
func (Homework) TableName() string {
	return "homeworks"
}

// FK制約の追加
func InitHomeworkFK() error {
	// HomeworkPosterUuid
	_, err := db.Exec("ALTER TABLE homeworks ADD FOREIGN KEY (homework_poster_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateHomeworkTestData() {
	afterSixMonth := time.Now().AddDate(0, 6, 0)
	afterOneYear := time.Now().AddDate(1, 0, 0)

	fmt.Println(afterSixMonth)
	fmt.Println(afterOneYear)

	homework1 := &Homework{
		HomeworkUuid:         "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		HomeworkLimit:        afterOneYear,
		TeachingMaterialUuid: "978f9835-5a16-4ac0-8581-7af8fac06b4e",
		StartPage:            24,
		PageCount:            2,
		HomeworkPosterUuid:   "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		HomeworkNote:         "がんばってくださ～い＾＾",
	}
	db.Insert(homework1)
	homework2 := &Homework{
		HomeworkUuid:         "K2079e71-3be5-4b4d-a0df-1f05859a7104",
		HomeworkLimit:        afterSixMonth,
		TeachingMaterialUuid: "978f9835-5a16-4ac0-8581-7af8fac06b4e",
		StartPage:            30,
		PageCount:            5,
		HomeworkPosterUuid:   "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		HomeworkNote:         "2こめ",
	}
	db.Insert(homework2)
	homework3 := &Homework{
		HomeworkUuid:         "K3079e71-3be5-4b4d-a0df-1f05859a7104",
		HomeworkLimit:        afterSixMonth,
		TeachingMaterialUuid: "978f9835-5a16-4ac0-8581-7af8fac06b4e",
		StartPage:            25,
		PageCount:            1,
		HomeworkPosterUuid:   "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		HomeworkNote:         "3こめ",
	}
	db.Insert(homework3)
}

// 宿題登録
func CreateHW(record Homework) error {
	affected, err := db.Insert(record)
	if err != nil { //エラーハンドル
		// XormのORMエラーを仕分ける
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // errをmysqlErrにアサーション出来たらtrue
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反
				return custom.NewErr(custom.ErrTypeUniqueConstraintViolation)
			default: // ORMエラーの仕分けにぬけがある可能性がある
				return custom.NewErr(custom.ErrTypeOtherErrorsInTheORM)
			}
		}
		// 通常の処理エラー
		return err
	}
	if affected == 0 {
		return custom.NewErr(custom.ErrTypeZeroEffectCUD)
	}

	return nil
}
