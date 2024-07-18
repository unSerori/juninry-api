package model

import (
	"time"
)

// おしらせテーブル
type Notice struct { // typeで型の定義, structは構造体
	NoticeUuid        string    `xorm:"varchar(36) pk" json:"noticeUUID"`        // おしらせの一意ID
	NoticeTitle       string    `xorm:"varchar(15) not null" json:"noticeTitle"` // おしらせのタイトル
	NoticeExplanatory string    `xorm:"text" json:"noticeExplanatory"`           // おしらせ内容
	NoticeDate        time.Time `xorm:"DATETIME not null" json:"noticeDate"`     // おしらせの時刻
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
		NoticeExplanatory: "来週の6/4(火)の遠足にて、おべんとうが必要です。また、同日にぞうきんの回収を行いますのでよろしくお願いします。,1",
		NoticeDate:        time.Now(),
		UserUuid:          "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		ClassUuid:         "09eba495-fe09-4f54-a856-9bea9536b661",
	}

	db.Insert(notice1)

	notice2 := &Notice{
		NoticeUuid:        "2097a7bb-5140-460d-807e-7173a51672bd",
		NoticeTitle:       "【持ち物】おべんと",
		NoticeExplanatory: "来週の6/4(火)の遠足にて、おべんとうが必要です。また、同日にぞうきんの回収を行いますので",
		NoticeDate:        time.Now(),
		UserUuid:          "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		ClassUuid:         "09eba495-fe09-4f54-a856-9bea9536b661",
	}

	db.Insert(notice2)

}

// 新規お知らせ登録
func CreateNotice(record Notice) (int64, error){
	
	affected, err := db.Insert(record)
	return affected, err
}

// classUuidで絞り込んだnoticeの結果を返す
func FindNotices(classUuids []string) ([]Notice, error) {

	// 結果を格納する変数宣言(findの結果)
	var notices []Notice

	//classUuidで絞り込んだデータを全件取得
	err := db.In("class_Uuid", classUuids).Find(
		&notices,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return nil, err
	}

	// エラーが出なければ取得結果を返す
	return notices, nil
}

// noticeUuidで絞り込んだnoticeの詳細を返す
func GetNoticeDetail(noticeUuid string) (Notice, error) {

	//結果格納用変数
	var noticeDetail Notice

	//noticeuuidで絞り込んで1件取得
	//.Getの返り値は存在の真偽値とエラー
	_, err := db.Where("notice_uuid =? ", noticeUuid).Get(
		&noticeDetail,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return Notice{}, err
	}

	return noticeDetail, nil
}

