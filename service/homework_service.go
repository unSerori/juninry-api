package service

import (
	"errors"
	"fmt"
	"juninry-api/model"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HomeworkService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// 課題データの構造体
type HomeworkData struct {
	HomeworkUuid              string `json:"homeworkUUID"` // 課題ID UUIDを大文字というきもち
	StartPage                 int    // 開始ページ
	PageCount                 int    // ページ数
	HomeworkNote              string // 課題の説明
	TeachingMaterialName      string // 教材名
	SubjectId                 int    // 教科ID
	SubjectName               string // 教科名
	TeachingMaterialImageUuid string `json:"TeachingMaterialImageUUID"` // 画像ID どういう扱いになるのかな UUIDを大文字というきもち
	ClassName                 string // クラス名
	SubmitFlag                int    // 提出フラグ 1 提出 0 未提出
}

// 締め切りごとに課題データをまとめた構造体
type TransformedData struct {
	HomeworkLimit time.Time      `json:"homeworkLimit"` //提出期限
	HomeworkData  []HomeworkData `json:"homeworkData"`  //課題データのスライス
}

// userUuidをuserHomeworkモデルに投げて、受け取ったデータを整形して返す
func (s *HomeworkService) FindHomework(userUuid string) ([]TransformedData, error) {

	//user_uuidを絞り込み条件にクソデカ構造体のスライスを受け取る
	userHomeworkList, err := model.FindUserHomework(userUuid)
	if err != nil { //エラーハンドル エラーを上に投げるだけ
		return nil, err
	}

	//期限をキー、バリューを課題データのマップにする
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

// 宿題登録処理
func (s *HomeworkService) SubmitHomework(c *gin.Context, bHW model.HomeworkSubmission, form *multipart.Form) error {
	// 画像の保存
	images := form.File["images"] // スライスからimages fieldを取得
	// 保存先ディレクトリ
	dst := "./upload/homework"
	// それぞれのファイルを保存
	for _, image := range images {
		fmt.Printf("image.Filename: %v\n", image.Filename)
		// ファイル名をuuidで作成
		fileName, err := uuid.NewRandom() // 新しいuuidの生成
		if err != nil {
			return err
		}
		// バリデーション
		// TODO: 形式(png, jpg, jpeg, gif, HEIF)
		// TODO: ファイルの種類->拡張子
		// TODO: パーミッション
		// 保存
		c.SaveUploadedFile(image, dst+"/"+fileName.String()+".png")
	}

	// 画像名スライスを文字列に変換し、
	// list :=
	// 画像一覧を提出中間テーブル構造体インスタンスに追加し、
	// bHW.list =
	// テーブルに追加。
	// ins

	return errors.New("hoge")
}
