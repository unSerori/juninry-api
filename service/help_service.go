package service

import (
	// "crypto/rand"
	"errors"
	"juninry-api/model"
	"time"

	"github.com/google/uuid"
)

type HelpService struct{}

// TODO:エラーハンドリング
// おてつだいを取得
func (s *HelpService) GetHelps(userUUID string) ([]model.Help, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return []model.Help{}, errors.New("user is not a teacher")
	}
	// ユーザー情報の取得
	bUser, err := model.GetUser(userUUID)
	if err != nil {
		return []model.Help{}, err
	}
	// ユーザー情報のouchiUUIDでご褒美を取得
	helps, err := model.GetHelps(*bUser.OuchiUuid, userUUID)
	if err != nil {
		return []model.Help{}, err
	}
	return helps, nil
}

// おてつだいを追加
func (s *HelpService) CreateHelp(userUUID string, help model.Help) (model.Help, error) {
	// ユーザーが教員であれば返す
	result, err := model.IsTeacher(userUUID)
	if result || err != nil {
		return model.Help{}, errors.New("user is not a teacher")
	}
	// 児童であっても同様
	result, err = model.IsJunior(userUUID)
	if result || err != nil {
		return model.Help{}, errors.New("user is not a junior")
	}

	// ごほうびUUIDの生成
	helpUUID, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {                   // 空の構造体とエラー
		return model.Help{}, err
	}
	help.HelpUuid = helpUUID.String() // uuidを文字列に変換してバインド

	// ごほうび作成
	// エラーが出なければコミットして追加したごほうびを返す
	_, err = model.CreateHelp(help)
	if err != nil {
		return model.Help{}, err
	}
	return help, err
}

// おてつだいを消化
func (s *HelpService) HelpDigestion(userUUID string, helpSubmittion model.HelpSubmittion) (*int, error) {
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

	helpSubmittion.UserUuid = userUUID       // ユーザーIDをバインド
	helpSubmittion.SubmittionAt = time.Now() // 現在時刻をバインド

	// おてつだいを消化
	// エラーが出なければコミットして追加したごほうびを返す
	_, err = model.StoreHelpSubmittion(helpSubmittion)
	if err != nil {
		return nil, err
	}
	// 交換できたらその分のポイントを減らす
	ouchiPoint, err := model.IncrementUpdatePoint(userUUID, helpSubmittion.HelpUuid)
	if err != nil {
		return nil, err
	}

	return ouchiPoint, err
}
