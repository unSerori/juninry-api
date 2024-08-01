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
	// ユーザー情報を取得
	user, err := model.GetUser(userUUID)
	if result || err != nil {
		return model.Reward{}, errors.New("user is not resorce")
	}
	reward.OuchiUuid = *user.OuchiUuid // おうちのuuidをバインド

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
func (s *RewardService) ExchangeReward(userUUID string, rewardExchange model.RewardExchanging) (*int, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return nil, errors.New("user is not a teacher")
		// 保護者であっても同様
	}
	result, err = model.IsPatron(userUUID)
	if result || err != nil {
		return nil, errors.New("user is not a teacher")
	}

	rewardExchange.ExchangingAt = time.Now() // 現在時刻をバインド
	rewardExchange.UserUuid = userUUID       // ユーザーIDをバインド

	// ごほうびを交換
	// エラーが出なければコミットして追加したごほうびを返す
	_, err = model.RewardExchange(rewardExchange)
	if err != nil {
		return nil, err
	}
	// 交換できたらその分のポイントを減らす
	point, err := model.DecrementUpdatePoint(userUUID, rewardExchange.RewardUuid)
	if err != nil {
		return nil, err
	}
	ouchiPoint := &point

	return ouchiPoint, err
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

// ごほうび交換履歴の構造体
type RewardExchangeHistory struct {
	RewardExchangengId int          `json:"rewardExchangengId"`
	ExchangengAt       string       `json:"exchangengAt"`
	UserName           string       `json:"userName"`
	RewardData         model.Reward `json:"rewardData"`
	IsExchange         int          `json:"isExchange"`
}

// 交換されたごほうびを取得
func (s *RewardService) GetRewardExchanging(userUUID string) ([]RewardExchangeHistory, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return []RewardExchangeHistory{}, errors.New("user is not a teacher")
	}
	// ouchiUUIDを取得
	user, err := model.GetUser(userUUID)
	if err != nil {
		return []RewardExchangeHistory{}, err
	}
	// 児童のIDを取得
	children, err := model.GetChildrenUuids(*user.OuchiUuid)
	if err != nil {
		return []RewardExchangeHistory{}, err
	}

	// ごほうび交換履歴の構造体
	var rewardExchangeHistories []RewardExchangeHistory
	// 交換されたごほうびを取得
	rewardExchangings, err := model.GetRewardExchangings(children)
	if err != nil {
		return []RewardExchangeHistory{}, err
	}
	// 交換されたごほうびを取得
	for _, rewardExchanging := range rewardExchangings {
		// ごほうびの取得
		reward, err := model.GetReward(rewardExchanging.RewardUuid)
		if err != nil {
			return []RewardExchangeHistory{}, err
		}
		// ユーザー名の取得
		user, err := model.GetUser(rewardExchanging.UserUuid)
		if err != nil {
			return []RewardExchangeHistory{}, err
		}
		// ごほうび交換履歴の構造体にバインド
		rewardExchangeHistory := RewardExchangeHistory{
			RewardExchangengId: rewardExchanging.RewardExchangingId,
			ExchangengAt:       rewardExchanging.ExchangingAt.Format("2006-01-02 15:04:05"),
			UserName:           user.UserName,
			RewardData:         reward,
			IsExchange:         rewardExchanging.Exchange,
		}
		// ごほうび交換履歴の構造体をスライスに追加
		rewardExchangeHistories = append(rewardExchangeHistories, rewardExchangeHistory)
	}
	return rewardExchangeHistories, nil
}
