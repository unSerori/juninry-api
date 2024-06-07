package service

import (
	"juninry-api/model"
)

type HomeworkService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// userUuidを引数として受け取り、model.Homeworkの配列を返すメソッド
func (s *HomeworkService) FindHomework(userUuid string) ([]model.UserHomework, error) {
	// //userUuidを元に所属構造体のスライスを取得する
	// classMemberships, err := model.FindClassMemberships(userUuid)
	// if err != nil {
	// 	return nil, err
	// }

	// //所属構造体のスライスからclassUuidのみ抽出
	// var classUuids []string
	// for _, classMembership := range classMemberships {
	// 	classUuids = append(classUuids, classMembership.ClassUuid)
	// }

	// //classUuidのスライスを元に教材構造体のスライスを取得する
	// teachingMaterials, err := model.FindTeachingMaterial(classUuids)
	// if err != nil {
	// 	return nil, err
	// }

	// //教材構造体のスライスからteachingMaterialUuidのみ抽出
	// var teachingMaterialUuids []string
	// for _, teachingMaterial := range teachingMaterials {
	// 	teachingMaterialUuids = append(teachingMaterialUuids, teachingMaterial.TeachingMaterialUuid)
	// }

	// //teachingMaterialUuidのスライスを元に課題構造体のスライスを取得する
	// homeworks, err := model.FindHomework(teachingMaterialUuids)

	//user_uuidを絞り込み条件にクソデカ構造体のスライスを受け取る
	homeworkList, err := model.FindUserHomework(userUuid)
	if err != nil { //エラーハンドル
		return nil, err
	}

	//できたら返す
	return homeworkList, nil
}
