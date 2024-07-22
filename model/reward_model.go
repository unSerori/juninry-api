package model

// 課題テーブル
type Reward struct { // typeで型の定義, structは構造体
	RewardUuid    string `xorm:"varchar(36) pk" json:"rewardUUID"`  // タスクのID
	OuchiUuid   string `xorm:"varchar(36) pk" json:"ouchiUUID"` // タスクのID
	RewardPoint int    `xorm:"not null" json:"rewardPoint"`     // 教材ID
	RewardNote    string `json:"rewardNote"`                        // 開始ページ
	RewardTitle   string `xorm:"not null" json:"rewardTitle"`       // ページ数
	IconId      int    `xorm:"not null" json:"iconId"`          // 投稿者ID
}

// テーブル名
func (Reward) TableName() string {
	return "rewards"
}

// FK制約の追加
func InitRewardFK() error {
	// HomeworkPosterUuid
	_, err := db.Exec("ALTER TABLE rewards ADD FOREIGN KEY (ouchi_uuid) REFERENCES ouchies(ouchi_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateRewardTestData() {
	help1 := &Reward{
		RewardUuid:    "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		OuchiUuid:   "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint: 10,
		RewardNote:    "200円まで",
		RewardTitle:   "アイス購入権",
		IconId:      1,
	}
	db.Insert(help1)
	help2 := &Reward{
		RewardUuid:    "a3579e71-3be5-4b4d-a0df-1f05859a7103",
		OuchiUuid:   "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint: 25,
		RewardNote:    "予算千円",
		RewardTitle:   "晩ごはん決定権",
		IconId:      2,
	}
	db.Insert(help2)
}

// 新規ごほうび登録
// 新しい構造体をレコードとして受け取り、ouchiテーブルにinsertし、成功した列数とerrorを返す
func CreateReward(record Reward) (int64, error) {
	affected, err := db.Nullable("invite_code", "valid_until").Insert(record)
	return affected, err
}

// 複数のごほうびを取得
func GetRewards(ouchiUuid string) ([]Reward, error) {
	//結果格納用変数
	var rewards []Reward
	//ouchiUuidで絞り込んで全取得
	err := db.Where("ouchi_uuid =?", ouchiUuid).Find(
		&rewards,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return []Reward{}, err
	}
	return rewards, nil
}

// TODO:ごほうびを更新

// ごほうびを削除
func DeleteReward(rewardUUID string) (int64, error) {
	//ouchiUuidで絞り込んで全取得
	result, err := db.Where("reward_uuid =?", rewardUUID).Delete()
	// データが取得できなかったらerrを返す
	if err != nil {
		return result, err
	}
	return result, nil
}
