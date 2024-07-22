package model

import (
	"time"
)

// ごほうび交換記録テーブル
type HelpSubmittion struct {
	HelpSubmittionId int       `xorm:"pk autoincr" json:"helpSubmittionId"` // おてつだいID
	UserUuid         string    `xorm:"varchar(36) pk" json:"userUUID"`      // ユーザーID
	HelpUuid         string    `xorm:"varchar(36) pk" json:"helpUUID"`      // おてつだいID
	SubmittionAt     time.Time `xorm:"TEXT" json:"submittionAt"`            // 提出日時
}

// テーブル名
func (HelpSubmittion) TableName() string {
	return "help_submittions"
}

// FK制約の追加
func InitHelpSubmittionFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE help_submittions ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// ClassUuid
	_, err = db.Exec("ALTER TABLE help_submittions ADD FOREIGN KEY (help_uuid) REFERENCES helps(help_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// おてつだいの記録を全件取得
func GetHelpSubmittions(userUuid string) ([]HelpSubmittion, error) {
	//結果格納用変数
	var helpSubmittions []HelpSubmittion
	err := db.Where("user_uuid = ?", userUuid).Find(&helpSubmittions)
	return helpSubmittions, err
}

// おてつだいの記録を登録
func StoreHelpSubmittion(help *HelpSubmittion) (bool, error) {
	affected, err := db.Insert(help)
	if err != nil || affected == 0 {
		return false, err
	}
	return true, err
}
