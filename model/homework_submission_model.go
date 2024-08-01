package model

// 宿題提出管理テーブル
type HomeworkSubmission struct {
	HomeworkUuid        string `xorm:"varchar(36) pk" json:"homeworkUUID" form:"homeworkUUID"` // ユーザーID
	UserUuid            string `xorm:"varchar(36) pk" json:"userUUID"`                         // クラスID
	ImageNameListString string `xorm:"TEXT" json:"imageNameListString"`                        // 画像ファイル名一覧 // TEXT型でUTF-8 21,845文字 // 一画像40文字と考えると最大546.125画像保存可能
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
