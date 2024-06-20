package service

import (
	"juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
)

type ClassService struct{}

func (s *ClassService) PermissionCheckedClassCreation(userUuid string, bClass model.Class) (model.Class, error) {

	// 返すやつを定義
	var class model.Class

	// クラス作成権限を持っているか確認
	isTeacher, err := model.CheckIsTeacher(userUuid)
	if err != nil { // エラーハンドル
		return class, err
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return class, common.NewErr(common.ErrTypePermissionDenied)
	}

	return class, nil
	// クラス作成処理

	// 招待コードの作成
	// 有効期限の設定

	// トランザクションしたい
	// _, err := model.CreateClass(bClass)
	// return

}
