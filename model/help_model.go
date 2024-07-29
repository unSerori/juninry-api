package model

import (
	"time"
)

// 課題テーブル
type Help struct { // typeで型の定義, structは構造体
	HelpUuid       string `xorm:"varchar(36) pk" json:"helpUUID"` // おてつだいのID
	OuchiUuid      string `xorm:"varchar(36)" json:"ouchiUUID"`   // おうちのID
	RewardPoint    int    `xorm:"not null" json:"rewardPoint"`    // もらえるおうちポイント
	HelpContent    string `json:"helpContent"`                    // 概要
	HelpTitle      string `xorm:"not null" json:"helpTitle"`      // タイトル
	IconId         int    `xorm:"not null" json:"iconId"`         // iconID
	SubmittedToday bool   `xorm:"-" json:"submittedToday"`        // フラグを追加　dbには保存されない
}

// テーブル名
func (Help) TableName() string {
	return "helps"
}

// FK制約の追加
func InitHelpFK() error {
	// HomeworkPosterUuid
	_, err := db.Exec("ALTER TABLE helps ADD FOREIGN KEY (ouchi_uuid) REFERENCES ouchies(ouchi_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateHelpTestData() {
	help1 := &Help{
		HelpUuid:    "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		OuchiUuid:   "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint: 24,
		HelpContent: "たたむところまでおねがいね",
		HelpTitle:   "せんたくもの",
		IconId:      1,
	}
	db.Insert(help1)
	help2 := &Help{
		HelpUuid:    "a3579e71-3be5-4b4d-a0df-1f05859a7103",
		OuchiUuid:   "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint: 25,
		HelpContent: "乾かしてるのは直してね",
		HelpTitle:   "あらいもの",
		IconId:      2,
	}
	db.Insert(help2)
}

// 新規おてつだい登録
// 新しい構造体をレコードとして受け取り、ouchiテーブルにinsertし、成功した列数とerrorを返す
func CreateHelp(record Help) (int64, error) {
	affected, err := db.Nullable("invite_code", "valid_until").Insert(record)
	return affected, err
}

// 特定のおてつだいを取得
func GetHelp(helpUUID string) (Help, error) {
	//結果格納用変数
	var help Help
	//ouchiUuidで絞り込んで全取得
	_, err := db.Where("help_uuid =?", helpUUID).Get(
		&help,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return Help{}, err
	}
	return help, nil
}

// 複数のおてつだいを取得
func GetHelps(ouchiUuid string, userUuid string) ([]Help, error) {
	//結果格納用変数
	var helps []Help
	today := time.Now().Format("2006-01-02")
	//ouchiUuidで絞り込んで全取得
	// err := db.Where("ouchi_uuid =?", ouchiUuid).Find(
	// 	&helps,
	// )
	// err := db.Table("helps").
	// 	Alias("h").
	// 	Join("LEFT", "help_submittions", "h.help_uuid = help_submittions.help_uuid AND DATE(help_submittions.submittion_at) = ? AND help_submittions.user_uuid = ?", today, userUuid).
	// 	Where("h.ouchi_uuid = ?", ouchiUuid).
	// 	Select("h.*, CASE WHEN help_submittions.help_submittion_id IS NOT NULL THEN 1 ELSE 0 END as submitted_today").
	// 	Find(&helps)
	err := db.Table("helps").
		Alias("h").
		Join("LEFT", "(SELECT help_uuid, MAX(submittion_at) as latest_submittion FROM help_submittions WHERE DATE(submittion_at) = ? AND user_uuid = ? GROUP BY help_uuid) hs", "h.help_uuid = hs.help_uuid", today, userUuid).
		Where("h.ouchi_uuid = ?", ouchiUuid).
		Select("h.*, CASE WHEN hs.latest_submittion IS NOT NULL THEN 1 ELSE 0 END as submitted_today").
		Find(&helps)
	// データが取得できなかったらerrを返す
	if err != nil {
		return []Help{}, err
	}
	return helps, nil
}

// TODO:おてつだいを更新

// おてつだいを削除
func DeleteHelp(helpUUID string) (int64, error) {
	//ouchiUuidで絞り込んで全取得
	result, err := db.Where("help_uuid =?", helpUUID).Delete()
	// データが取得できなかったらerrを返す
	if err != nil {
		return result, err
	}
	return result, nil
}
