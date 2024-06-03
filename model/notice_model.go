package model

import (
	"time"
)

// おしらせテーブル
type Notice struct { // typeで型の定義, structは構造体
	NoticeUuid        string    `xorm:"varchar(36) pk" json:"noticeUUID"`        // おしらせの一意ID
	NoticeTitle       string    `xorm:"varchar(15) not null" json:"noticeTitle"` // おしらせのタイトル
	NoticeExplanatory string    `xorm:"text" json:"noticeExplanatory"`           // おしらせ内容
	NoticeDate        time.Time `xorm:"DATE not null" json:"noticeDate"`         // おしらせの時刻
	UserUuid          string    `xorm:"varchar(36) not null" json:"userUUID"`    // おしらせ発行ユーザ
	ClassUuid         string    `xorm:"varchar(36) not null" json:"classUUID"`   // どのクラスのお知らせか
}

// テーブル名
func (Notice) TableName() string {
	return "notices"
}

// FK制約の追加
func InitNoticeFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE notices ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// ClassUuid
	_, err = db.Exec("ALTER TABLE notices ADD FOREIGN KEY (class_uuid) REFERENCES classes(class_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateNoticeTestData() {
	notice1 := &Notice{
		NoticeUuid:        "51e6807b-9528-4a4b-bbe2-d59e9118a70d",
		NoticeTitle:       "【持ち物】おべんとうとぞうきん",
		NoticeExplanatory: "来週の6/4(火)の遠足にて、おべんとうが必要です。また、同日にぞうきんの回収を行いますのでよろしくお願いします。",
		NoticeDate:        time.Now(),
		UserUuid:          "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		ClassUuid:         "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(notice1)
}
