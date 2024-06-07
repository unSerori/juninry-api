package model

// ユーザのクラス所属中間テーブル
type NoticeReadStatus struct {
	NoticeUuid string `xorm:"varchar(36) pk" json:"noticeUUID"` // おしらせID
	UserUuid   string `xorm:"varchar(36) pk" json:"userUUID"`   // ユーザーID
}

// テーブル名
func (NoticeReadStatus) TableName() string {
	return "notice_read_statuses"
}

// FK制約の追加
func InitNoticeReadStatus() error {
	// NoticeUuid
	_, err := db.Exec("ALTER TABLE notice_read_statuses ADD FOREIGN KEY (notice_uuid) REFERENCES notices(notice_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// UserUuid
	_, err = db.Exec("ALTER TABLE notice_read_statuses ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateNoticeReadStatusTestData() {
	nrs1 := &NoticeReadStatus{
		NoticeUuid: "51e6807b-9528-4a4b-bbe2-d59e9118a70d",
		UserUuid:   "3cac1684-c1e0-47ae-92fd-6d7959759224",
	}
	db.Insert(nrs1)
}