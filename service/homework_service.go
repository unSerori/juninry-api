package service

import (
	"juninry-api/model"
	"time"
)

type HomeworkService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// 課題データの構造体
type HomeworkData struct {
	HomeworkUuid              string `json:"homework_uuid"`                // 課題ID
	StartPage                 int    `json:"start_page"`                   // 開始ページ
	PageCount                 int    `json:"page_count"`                   // ページ数
	HomeworkNote              string `json:"homework_note"`                // 課題の説明
	TeachingMaterialName      string `json:"teaching_material_name"`       // 教材名
	SubjectId                 int    `json:"subject_id"`                   // 教科ID
	SubjectName               string `json:"subject_name"`                 // 教科名
	TeachingMaterialImageUuid string `json:"teaching_material_image_uuid"` // 画像ID どういう扱いになるのかな
	ClassName                 string `json:"class_name"`                   // クラス名
	SubmitFlag                int    `json:"submit_flag"`                  // 提出フラグ 1 提出 0 未提出
}

// 締め切りごとに課題データをまとめた構造体
type TransformedData struct {
	HomeworkLimit time.Time      `json:"homework_limit"` //提出期限
	HomeworkData  []HomeworkData `json:"homework_data"`  //課題データのスライス
}

// userUuidをuserHomeworkモデルに投げて、受け取ったデータを整形して返す
func (s *HomeworkService) FindHomework(userUuid string) ([]TransformedData, error) {

	//user_uuidを絞り込み条件にクソデカ構造体のスライスを受け取る
	userHomeworkList, err := model.FindUserHomework(userUuid)
	if err != nil { //エラーハンドル エラーを上に投げるだけ
		return nil, err
	}

	//一旦期限をキー、バリューを課題データのマップにする
	transformedDataMap := make(map[time.Time][]HomeworkData)
	for _, userHomework := range userHomeworkList {
		homeworkData := HomeworkData{
			HomeworkUuid:              userHomework.HomeworkUuid,
			StartPage:                 userHomework.StartPage,
			PageCount:                 userHomework.PageCount,
			HomeworkNote:              userHomework.HomeworkNote,
			TeachingMaterialName:      userHomework.TeachingMaterialName,
			SubjectId:                 userHomework.SubjectId,
			SubjectName:               userHomework.SubjectName,
			TeachingMaterialImageUuid: userHomework.TeachingMaterialImageUuid,
			ClassName:                 userHomework.ClassName,
			SubmitFlag:                userHomework.SubmitFlag,
		}
		transformedDataMap[userHomework.HomeworkLimit] = append(transformedDataMap[userHomework.HomeworkLimit], homeworkData)
	}

	//作ったマップをさらに整形
	var transformedDataList []TransformedData
	for limit, homeworkData := range transformedDataMap {
		transformedData := TransformedData{
			HomeworkLimit: limit,
			HomeworkData:  homeworkData,
		}
		transformedDataList = append(transformedDataList, transformedData)
	}

	//できたら返す
	return transformedDataList, nil
}
