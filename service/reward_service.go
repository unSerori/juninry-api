package service

import (
	"errors"
	"juninry-api/model"
	"time"

	"github.com/google/uuid"
)

type RewardService struct{}

// ごほうびを一件取得
func (s *RewardService) GetReward(rewardUUID string) (model.Reward, error) {
	// ユーザー情報のouchiUUIDでご褒美を取得
	reward, err := model.GetReward(rewardUUID)
	if err != nil {
		return model.Reward{}, err
	}
	return reward, nil
}

// TODO:エラーハンドリング
// ごほうびを取得
func (s *RewardService) GetRewards(userUUID string) ([]model.Reward, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return []model.Reward{}, errors.New("user is not a teacher")
	}
	// ユーザー情報の取得
	bUser, err := model.GetUser(userUUID)
	if err != nil {
		return []model.Reward{}, err
	}
	// ユーザー情報のouchiUUIDでご褒美を取得
	rewards, err := model.GetRewards(*bUser.OuchiUuid)
	if err != nil {
		return []model.Reward{}, err
	}
	return rewards, nil
}

// ごほうびを追加
func (s *RewardService) CreateRewards(userUUID string, reward model.Reward) (model.Reward, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return model.Reward{}, errors.New("user is not a teacher")
	}
	// 児童であっても同様
	result, err = model.IsJunior(userUUID)
	if result || err != nil {
		return model.Reward{}, errors.New("user is not a teacher")
	}

	// ごほうびUUIDの生成
	rewardUuid, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {                     // 空の構造体とエラー
		return model.Reward{}, err
	}
	reward.RewardUuid = rewardUuid.String() // uuidを文字列に変換してバインド

	// ごほうび作成
	// エラーが出なければコミットして追加したごほうびを返す
	_, err = model.CreateReward(reward)
	if err != nil {
		return model.Reward{}, err
	}
	return reward, err
}

// ごほうびを削除
func (s *RewardService) DeleteReward(userUUID string, reward_UUID string) (bool, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return false, errors.New("user is not a teacher")
	}
	// 児童であっても同様
	result, err = model.IsJunior(userUUID)
	if result || err != nil {
		return false, errors.New("user is not a teacher")
	}

	// ごほうびを削除
	// エラーが出なければコミットして追加したごほうびを返す
	_, err = model.DeleteReward(reward_UUID)
	if err != nil {
		return false, err
	}
	return true, err
}

// ごほうびを交換
func (s *RewardService) ExchangeReward(userUUID string, rewardExchange model.RewardExchanging) (bool, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return false, errors.New("user is not a teacher")
		// 保護者であっても同様
	}
	result, err = model.IsPatron(userUUID)
	if result || err != nil {
		return false, errors.New("user is not a teacher")
	}

	rewardExchange.ExchangingAt = time.Now() // 現在時刻をバインド

	// ごほうびを交換
	// エラーが出なければコミットして追加したごほうびを返す
	dResult, err := model.RewardExchange(rewardExchange)
	if err != nil {
		return false, err
	}
	// 交換できたらその分のポイントを減らす
	_, err = model.DecrementUpdatePoint(userUUID, rewardExchange.RewardUuid)

	return dResult, err
}

// 交換されたごほうびを消化
func (s *RewardService) RewardDigestion(userUUID string, rewardExchangingId string) (bool, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return false, errors.New("user is not a teacher")
	}
	// juniorであっても同様
	result, err = model.IsJunior(userUUID)
	if result || err != nil {
		return false, errors.New("user is not a teacher")
	}
	// ごほうびを交換
	// エラーが出なければコミットして追加したごほうびを返す
	dResult, err := model.UpdateRewardExchanging(rewardExchangingId)
	if err != nil {
		return false, err
	}
	return dResult, err
}
