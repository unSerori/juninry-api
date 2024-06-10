package service

import (
	"juninry-api/model"
)

type NoticeService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// こめんとかいて～～
func (s *NoticeService) FindNotices(userUuid string) ([]model.Notice, error) {

	// userUuidを条件にしてclassUuidを取ってくる
	// 1 - userUuidからclass_menbershipの構造体を取ってくる
	classMemberships, err := model.FindClassMemberships(userUuid)
	if err != nil {
		return nil, err
	}

	// 2 - 構造体からclassUuidのスライス(配列)を作る
	var classUuids []string
	for _, classMembership := range classMemberships {
		classUuid := classMembership.ClassUuid
		classUuids = append(classUuids, classUuid)
	}

	// classUuidを条件にしてnoticeの構造体を取ってくる
	notices, err := model.FindNotices(classUuids)
	if err != nil {
		return nil, err
	}

	return notices, nil
}
