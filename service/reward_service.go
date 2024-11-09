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
	// RewardName   string `json:"rewardName"`
	// RewardPoint  int    `json:"rewardPoint"`
	// RewardTitle  string `json:"rewardTitle"`
	// IconId       int    `json:"iconId"`
	DepositPoint int           `json:"depositPoint"`
	BoxStatus    int           `json:"boxStatus"`
	Reward       *model.Reward `json:"reward"`
}

func (s *RewardService) GetBoxRewards(userUUID string) ([]BoxReward, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return nil, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// ユーザーが児童であればレスポンスの種類を制限する
	isJunior, err := model.IsJunior(userUUID)
	if err != nil {
		return nil, err
	}

	// ユーザー情報の取得
	bUser, err := model.GetUser(userUUID)
	if err != nil {
		return nil, err
	}

	// ユーザー情報のouchiUUIDでご褒美を取得
	boxes, err := model.GetBoxes(*bUser.OuchiUuid)
	if err != nil {
		return nil, err
	}

	// レスポンスの形式変換
	var boxRewards []BoxReward

	for _, box := range boxes {
		// レスポンスの形式変換
		boxReward := BoxReward{
			HardwareUuid: box.HardwareUuid,
			DepositPoint: box.DepositPoint,
			BoxStatus:    box.BoxStatus,
		}

		// 児童であればステータスが1でも2でもない場合はスキップ
		if isJunior && box.BoxStatus != 1 && box.BoxStatus != 2 {
			continue
		}

		// ボックスに対してご褒美が設定されているかを確認
		exist, err := model.BoxRewardExists(box.HardwareUuid)
		if err != nil {
			return nil, err
		}
		if exist {
			// ボックスのごほうびを取得
			reward, err := model.GetBoxReward(*bUser.OuchiUuid, box.HardwareUuid)
			if err != nil {
				return nil, err
			}

			boxReward.Reward = &reward
		}

		// レスポンスの形式変換
		boxRewards = append(boxRewards, boxReward)

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

	// 宝箱と紐づいたご褒美であればボックスのステータスを1に変更
	if reward.HardwareUuid != nil {
		model.UpdateBoxStatus(*reward.HardwareUuid, 1)
	}

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

	// ポイントの下限を確認
	if addPoint <= 0 {
		return 0, custom.NewErr(custom.ErrTypeUnforeseenCircumstances)
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

	// ボックスのポイントを更新
	err = model.UpdateBoxDepositPoint(hardUuid, boxCurrentPoint+addPoint)
	if err != nil {
		return 0, err
	}

	// ポイントがマックスであれば、ボックスの状態を変更する
	if boxCurrentPoint+addPoint == boxMaxPoint {
		err = model.UpdateBoxStatus(hardUuid, 2)
		if err != nil {
			return 0, err
		}
	}

	// ユーザーのポイントを更新
	err = model.UpdateOuchiPoint(userUUID, havePoint-addPoint)
	if err != nil {
		return 0, err
	}

	// 更新が完了したらtrueを返す
	return boxCurrentPoint + addPoint, nil
}

// ボックスのロック状態を入れ替える
func (s *RewardService) ChangeBoxLockStatus(userUUID string, hardUuid string) (int, error) {
	// ユーザーが保護者以外であれば返す
	result, err := model.IsPatron(userUUID)
	if err != nil {
		return 0, err
	}
	if !result {
		return 0, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// 宝箱がおうちのものであることを確認
	// ユーザーの詳細を取得
	user, err := model.GetUser(userUUID)
	if err != nil {
		return 0, err
	}
	if user.OuchiUuid == nil { // おうちのuuidを取得できなければ400
		return 0, custom.NewErr(custom.ErrTypeNoResourceExist)
	}
	ouchiUuid := *user.OuchiUuid // おうちのuuid

	// ボックスの詳細を取得
	box, err := model.GetBox(hardUuid)
	if err != nil {
		return 0, err
	}

	// ボックスがおうちのものであることを確認
	if box.OuchiUuid != ouchiUuid { // おうちのuuidを取得できなければ400
		return 0, custom.NewErr(custom.ErrTypeNoResourceExist)
	}


	var updateStatusNum int
	// ボックスのロック状態を入れ替える
	if box.BoxStatus == 1 || box.BoxStatus == 2 {
		updateStatusNum = 3
	} else if box.BoxStatus == 3 {
		result, err := model.BoxRewardExists(hardUuid)
		if err != nil {
			return 0, err
		}
		if result {
			reward, err := model.GetBoxReward(ouchiUuid, hardUuid)
			if err != nil {
				return 0, err
			}
			if reward.RewardPoint > box.DepositPoint {
				updateStatusNum = 1
			} else {
				updateStatusNum = 2
			}
		}
	}
	fmt.Println("box.BoxStatus")
	fmt.Println(box.BoxStatus)
	fmt.Println(updateStatusNum)
	// ボックスのロック状態を入れ替える
	err = model.UpdateBoxStatus(hardUuid, updateStatusNum)
	if err != nil {
		fmt.Println("UpdateBoxStatus error")
		fmt.Println(err)
		return 0, err
	}

	return updateStatusNum, nil
}
