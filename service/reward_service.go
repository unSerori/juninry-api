package service

import (
	"errors"
	"fmt"
	"juninry-api/common/logging"
	"juninry-api/model"
	"juninry-api/utility/custom"
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

type BoxReward struct {
	HardwareUuid string `json:"hardwareUUID"`
	RewardName   string `json:"rewardName"`
	RewardPoint  int    `json:"rewardPoint"`
	RewardTitle  string `json:"rewardTitle"`
	IconId       int    `json:"iconId"`
	DepositPoint int    `json:"depositPoint"`
}

func (s *RewardService) GetBoxRewards(userUUID string) ([]BoxReward, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return nil, errors.New("user is not a teacher")
	}
	// ユーザー情報の取得
	bUser, err := model.GetUser(userUUID)
	if err != nil {
		return nil, err
	}
	// ユーザー情報のouchiUUIDでご褒美を取得
	rewards, err := model.GetBoxRewards(*bUser.OuchiUuid)
	if err != nil {
		return nil, err
	}
	fmt.Println(rewards)

	// レスポンスの形式変換
	var boxRewards []BoxReward

	for _, reward := range rewards {
		// ハードウェアUUIDから現在のポイントを取得
		depositPoint, err := model.GetBoxDepositPoint(*reward.HardwareUuid)
		if err != nil {
			return nil, err
		}

		boxRewards = append(boxRewards, BoxReward{
			HardwareUuid: *reward.HardwareUuid,
			RewardName:   reward.RewardTitle,
			RewardPoint:  reward.RewardPoint,
			RewardTitle:  reward.RewardTitle,
			IconId:       reward.IconId,
			DepositPoint: depositPoint,
		})

	}
	return boxRewards, nil
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



// 宝箱にポイントを追加
func (s *RewardService) BoxAddPoint(userUUID string, addPoint int, hardUuid string) (int, error) {
	// ユーザーが子供であることを確認
	result, err := model.IsJunior(userUUID)
	if err != nil {
		return 0, err
	}
	if !result { // 子供以外は403
		return 0, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// 宝箱が自身のものであることを確認
	// ユーザーの詳細を取得
	user, err := model.GetUser(userUUID)
	if err != nil {
		return 0, err
	}
	if user.OuchiUuid == nil { // おうちのuuidを取得できなければ400
		return 0, custom.NewErr(custom.ErrTypeNoResourceExist)
	}
	ouchiUuid := *user.OuchiUuid // おうちのuuid

	// 現在のポイントを取得
	havePoint := user.OuchiPoint
	if havePoint < addPoint { // ポイント不足
		return 0, custom.NewErr(custom.ErrTypeUnforeseenCircumstances)
	}

	// 自身の所有する箱のごほうびを取得
	reward, err := model.GetBoxReward(ouchiUuid, hardUuid)
	if err != nil {
		return 0, err
	}
	if reward.RewardPoint == 0 { // 結果がなければ403
		logging.ErrorLog("Could not find the reward.", nil)
		return 0, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	boxMaxPoint := reward.RewardPoint // 宝箱に必要なポイント

	// ボックスの現在のポイントを取得
	boxCurrentPoint, err := model.GetBoxDepositPoint(hardUuid)
	if err != nil {
		return 0, err
	}

	// ポイントを追加しすぎてないかを確認
	if boxCurrentPoint+addPoint > boxMaxPoint { // ボックスのポイント上限を超えてる
		return 0, custom.NewErr(custom.ErrTypeUnforeseenCircumstances)
	}

	fmt.Println("宝箱にポイントを追加しました")
	fmt.Println("宝箱のポイント:", boxCurrentPoint)
	fmt.Println("追加するポイント:", addPoint)
	fmt.Println("宝箱のポイント上限:", boxMaxPoint)

	// ボックスのポイントを更新
	err = model.UpdateBoxDepositPoint(hardUuid, boxCurrentPoint+addPoint)
	if err != nil {
		fmt.Println(err)
		fmt.Println("ボックスのポイントを更新に失敗しました")
		return 0, err
	}

	// ユーザーのポイントを更新
	err = model.UpdateOuchiPoint(userUUID, havePoint-addPoint)
	if err != nil {
		fmt.Println(err)
		fmt.Println("ユーザーのポイントを更新に失敗しました")
		return 0, err
	}

	
	// 更新が完了したらtrueを返す
	return boxCurrentPoint + addPoint, nil
}
