package model

// 宿題提出管理テーブル
type HomeworkSubmission struct {
	HomeworkUuid string `xorm:"varchar(36) pk" json:"homeworkUUID"` // ユーザーID
	UserUuid     string `xorm:"varchar(36) pk" json:"userUUID"`     // クラスID
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
		HomeworkUuid: "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		UserUuid:     "3cac1684-c1e0-47ae-92fd-6d7959759224",
	}
	db.Insert(hs1)
}
