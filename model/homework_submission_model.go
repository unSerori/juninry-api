package model

import (
	"juninry-api/utility/custom"
	"time"
)

// 宿題提出管理テーブル
type HomeworkSubmission struct {
	HomeworkUuid        string    `xorm:"varchar(36) pk" json:"homeworkUUID" form:"homeworkUUID"` // ユーザーID
	UserUuid            string    `xorm:"varchar(36) pk" json:"userUUID"`                         // クラスID
	ImageNameListString string    `xorm:"TEXT" json:"imageNameListString"`                        // 画像ファイル名一覧 // TEXT型でUTF-8 21,845文字 // 一画像40文字と考えると最大546.125画像保存可能
	SubmissionDate      time.Time `xorm:"DATETIME not null" json:"submissionDate"`                // 提出日時
}

// テーブル名
func (HomeworkSubmission) TableName() string {
	return "homework_submissions"
}

// FK制約の追加
func InitHomeworkSubmissionFK() error {
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
func CreateHomeworkSubmissionTestData() {
	hs1 := &HomeworkSubmission{
		HomeworkUuid:        "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		UserUuid:            "3cac1684-c1e0-47ae-92fd-6d7959759224",
		ImageNameListString: "bbbbbbbb-a6ad-4059-809c-6df866e7c5e6.jpg, gggggggg-176f-4dea-bec0-21464f192869.jpg, rrrrrrrr-bb84-4565-9666-d53dfcb59dd3.jpg",
	}
	db.Insert(hs1)
}

// 提出構造体を登録
func StoreHomework(hwS *HomeworkSubmission) (bool, error) {
	affected, err := db.Insert(hwS)
	if err != nil || affected == 0 {
		return false, err
	}
	return true, err
}

// 課題提出状況の確認
func CheckHomeworkSubmission(homeworkUuids []string) (int64, error) {
	count, err := db.In("homework_uuid", homeworkUuids).Count(&HomeworkSubmission{})

	if err != nil {
		return -1, err
	}
	return count, nil
}

// 提出状況を取得
func GetHwSubmission(hwId string, userId string) (HomeworkSubmission, error) {
	var hwS HomeworkSubmission // 取得したデータをマッピングする構造体
	isFound, err := db.
		Where("homework_uuid = ?", hwId).
		Where("user_uuid = ?", userId).Get(&hwS)
	if err != nil {
		return HomeworkSubmission{}, err
	}
	if !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`!isFound`で、影響を与えるSQL文の時は`affected == 0`でハンドリング
		return HomeworkSubmission{}, custom.NewErr(custom.ErrTypeNoFoundR)
	}

	return hwS, nil
}
