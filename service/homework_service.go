package service

import (
	"juninry-api/model"
)

type HomeworkService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// userUuidを引数として受け取り、model.Homeworkの配列を返すメソッド
func (s *HomeworkService) FindHomework(userUuid string) ([]model.Homework, error) {
	//userUuidを元に所属クラスの構造体を取得する
	classMemberships, err := model.FindClassMemberships(userUuid)
	if err != nil {
		return nil, err
	}

	//所属構造体からクラスIDを配列に整形
	var classUuids []string
	for _, classMembership := range classMemberships {
		classUuids = append(classUuids, classMembership.ClassUuid)
	}

	//クラスIDの配列を元に教材の構造体を取得する
	teachingMaterials, err := model.FindTeachingMaterial(classUuids)
	if err != nil {
		return nil, err
	}

	//教材構造体から教材IDを配列に整形
	var teachingMaterialUuids []string
	for _, teachingMaterial := range teachingMaterials {
		teachingMaterialUuids = append(teachingMaterialUuids, teachingMaterial.TeachingMaterialUuid)
	}

	//教材IDの配列を元に課題一覧の構造体を取得する
	homeworks, err := model.FindHomework(teachingMaterialUuids)
	if err != nil {
		return nil, err
	}

	//できたら返す
	return homeworks, nil
}
