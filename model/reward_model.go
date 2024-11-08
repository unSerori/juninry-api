package model

// ご褒美テーブル
type Reward struct { // typeで型の定義, structは構造体
	RewardUuid    string  `xorm:"varchar(36) pk" json:"rewardUUID"`        // タスクのID
	OuchiUuid     string  `xorm:"varchar(36)" json:"ouchiUUID"`            // タスクのID
	RewardPoint   int     `xorm:"not null" json:"rewardPoint"`             // 教材ID
	RewardContent string  `xorm:"varchar(20)" json:"rewardContent"`        // 開始ページ
	RewardTitle   string  `xorm:"varchar(10) not null" json:"rewardTitle"` // ページ数
	IconId        int     `xorm:"not null" json:"iconId"`                  // 投稿者ID
	HardwareUuid  *string `xorm:"varchar(36)" json:"hardwareUUID"`         // ハードウェアUUID
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

	_, err = db.Exec("ALTER TABLE rewards ADD FOREIGN KEY (hardware_uuid) REFERENCES hardwares(hardware_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateRewardTestData() {
	help1 := &Reward{
		RewardUuid:    "a3579e71-3be5-4b4d-a0df-1f05859a7104",
		OuchiUuid:     "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint:   10,
		RewardContent: "200円まで",
		RewardTitle:   "アイス購入権",
		IconId:        1,
	}
	db.Insert(help1)
	help2 := &Reward{
		RewardUuid:    "a3579e71-3be5-4b4d-a0df-1f05859a7103",
		OuchiUuid:     "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint:   25,
		RewardContent: "予算千円",
		RewardTitle:   "晩ごはん決定権",
		IconId:        2,
	}
	db.Insert(help2)

	help3 := &Reward{
		RewardUuid:    "90b29b3f-190c-4fd7-a968-a5cd68086075",
		OuchiUuid:     "2e17a448-985b-421d-9b9f-62e5a4f28c49",
		RewardPoint:   1000,
		RewardContent: "Switchのゲームなんでも買います",
		RewardTitle:   "ゲーム交換券",
		IconId:        1,
	}
	hardUuid := "df2b1f4c-b49a-4068-80c5-3120dceb14c8"

	help3.HardwareUuid = &hardUuid
	db.Insert(help3)

}

// 新規ごほうび登録
// 新しい構造体をレコードとして受け取り、ouchiテーブルにinsertし、成功した列数とerrorを返す
func CreateReward(record Reward) (int64, error) {
	affected, err := db.Nullable("invite_code", "valid_until").Insert(record)
	return affected, err
}

// 複数のごほうびを取得
func GetReward(rewardUUID string) (Reward, error) {
	//結果格納用変数
	var reward Reward
	//ouchiUuidで絞り込んで全取得
	_, err := db.Where("reward_uuid =?", rewardUUID).Get(
		&reward,
	)
	// データが取得できなかったらerrを返す
	if err != nil {
		return Reward{}, err
	}
	return reward, nil
}

// 複数のごほうびを取得
func GetRewards(ouchiUuid string) ([]Reward, error) {
	//結果格納用変数
	var rewards []Reward
	//ouchiUuidで絞り込んで全取得
	err := db.Where("ouchi_uuid =? and hardware_uuid is null", ouchiUuid).Find(
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

// 自身の所有する箱のごほうびを取得
func GetBoxReward(ouchiUuid string, hardwareUuid string) (Reward, error) {
	// 結果格納用変数
	var reward Reward

	// ouchiUuidとhardwareUuidで絞り込んで一致するレコードを取得
	_, err := db.Where("ouchi_uuid =? and hardware_uuid =?", ouchiUuid, hardwareUuid).Get(&reward)
	if err != nil {
		return Reward{}, err
	}

	return reward, nil
}

//
