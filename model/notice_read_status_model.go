package model

// ユーザのクラス所属中間テーブル
type NoticeReadStatus struct {
	NoticeUuid string `xorm:"varchar(36) pk" json:"noticeUUID"` // おしらせID
	OuchiUuid  string `xorm:"varchar(36) pk" json:"ouchiUUID"`  // おうちID
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
	_, err = db.Exec("ALTER TABLE notice_read_statuses ADD FOREIGN KEY (ouchi_uuid) REFERENCES ouchies(ouchi_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateNoticeReadStatusTestData() {
	nrs1 := &NoticeReadStatus{
		NoticeUuid: "51e6807b-9528-4a4b-bbe2-d59e9118a70d",
		OuchiUuid:  "2e17a448-985b-421d-9b9f-62e5a4f28c49",
	}
	db.Insert(nrs1)
}

// notice_read_statusにデータがあるか調べる(確認済みの場合、データが存在する)
func IsRead(noticeUuid string, ouchiUuid string) (bool, error) {

	//noticeUuidとuserUuidから一致するデータがあるか取得
	has, err := db.Where("notice_uuid = ? AND ouchi_uuid = ?", noticeUuid, ouchiUuid).Exist(&NoticeReadStatus{})
	if err != nil {
		return false, err
	}

	return has, nil
}

// noticeUuidで絞った結果を返す
func GetNoticeStatusList(noticeUuid string) ([]NoticeReadStatus, error) {

	// 結果を格納する変数宣言(findの結果)
	var noticeReadStatus []NoticeReadStatus

	// noticeUuidで条件指定
	err := db.Where("notice_uuid = ?", noticeUuid).Find(&noticeReadStatus)
	// データが取得できなかったらerrを返す
	if err != nil {
		return nil, err
	}

	// エラーが出なければ取得結果を返す
	return noticeReadStatus, nil
}

// お知らせ確認登録
func ReadNotice(record NoticeReadStatus) error {
	attached, err := db.Insert(record)
	if err != nil || attached == 0 {
		return err
	}
	return err
}
