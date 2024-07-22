package model

import (
	"time"
)

// ごほうび交換記録テーブル
type RewardExchanging struct {
	RewardExchangingId int       `xorm:"pk autoincr" json:"rewardExchangingId"`       // 交換記録ID
	RewardUuid   string    `xorm:"varchar(36) pk" json:"rewardUUID" form:"rewardUUID"` // ユーザーID
	UserUuid     string    `xorm:"varchar(36) pk" json:"userUUID"`                     // クラスID
	ExchangingAt time.Time `xorm:"TEXT" json:"exchangingAt"`                           // 画像ファイル名一覧 // TEXT型でUTF-8 21,845文字 // 一画像40文字と考えると最大546.125画像保存可能
	Exchange     int       `xorm:"int" json:"exchange"`                                // 0 まだ or 1 交換済
}

// テーブル名
func (RewardExchanging) TableName() string {
	return "reward_exchangings"
}

// FK制約の追加
func InitRewardExchangingFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE reward_exchangings ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	// ClassUuid
	_, err = db.Exec("ALTER TABLE reward_exchangings ADD FOREIGN KEY (reward_uuid) REFERENCES rewards(reward_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// 交換記録を全件取得
func GetRewardExchangings(userUuid string) ([]RewardExchanging, error) {
	//結果格納用変数
	var rewardExchangings []RewardExchanging
	err := db.Where("user_uuid = ?", userUuid).Find(&rewardExchangings)
	return rewardExchangings, err
}

// 交換記録を登録
func StoreRewardExchanging(reward *RewardExchanging) (bool, error) {
	affected, err := db.Insert(reward)
	if err != nil || affected == 0 {
		return false, err
	}
	return true, err
}

// 交換記録の処理を更新　０→１に
func UpdateRewardExchanging(reward *RewardExchanging) (bool, error) {
	affected, err := db.Where("reward_exchanging_id = ?", reward.RewardExchangingId).Update(&RewardExchanging{Exchange: 1})
	if err != nil || affected == 0 {
		return false, err
	}
	return true, err
}
