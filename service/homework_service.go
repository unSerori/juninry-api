package service

import (
	"errors"
	"fmt"
	"juninry-api/common"
	"juninry-api/dip"
	"juninry-api/model"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type HomeworkService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// 課題データの構造体
type HomeworkData struct {
	HomeworkUuid              string `json:"homeworkUUID"`              // 課題ID
	StartPage                 int    `json:"startPage"`                 // 開始ページ
	PageCount                 int    `json:"pageCount"`                 // ページ数
	HomeworkNote              string `json:"homeworkNote"`              // 課題の説明
	TeachingMaterialName      string `json:"teachingMaterialName"`      // 教材名
	SubjectId                 int    `json:"subjectId"`                 // 教科ID
	SubjectName               string `json:"subjectName"`               // 教科名
	TeachingMaterialImageUuid string `json:"teachingMaterialImageUUID"` // 画像ID どういう扱いになるのかな
	ClassName                 string `json:"className"`                 // クラス名
	SubmitFlag                int    `json:"submitFlag"`                // 提出フラグ 1 提出 0 未提出
}

// 締め切りごとに課題データをまとめた構造体
type TransformedData struct {
	HomeworkLimit time.Time      `json:"homeworkLimit"` //提出期限
	HomeworkData  []HomeworkData `json:"homeworkData"`  //課題データのスライス
}

// userUuidをuserHomeworkモデルに投げて、受け取ったデータを整形して返す
func (s *HomeworkService) FindHomework(userUuid string) ([]TransformedData, error) {

	// 親には宿題一覧使えないよ
	isPatron, err := model.IsPatron(userUuid)
	if err != nil {
		return nil, err
	}
	if isPatron {	// 親が宿題一覧見ようとしないでね、何も情報とれないんだけどさ、、、
		return nil, common.NewErr(common.ErrTypePermissionDenied)
	}

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
// インターフェース型で依存性を受け取ることにより、具体的な実装(gin.Context, GinContextWrapper)ではなくインターフェースに依存し、依存性逆転が実現できる。
func (s *HomeworkService) SubmitHomework(uploader dip.FileUpLoader, bHW model.HomeworkSubmission, form *multipart.Form) error {
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
		uploader.SaveUploadedFile(image, dst+"/"+fileName.String()+".png") // c.SaveUploadedFile(image, dst+"/"+fileName.String()+".png")
	}

	// 画像名スライスを文字列に変換し、
	// list :=
	// 画像一覧を提出中間テーブル構造体インスタンスに追加し、
	// bHW.list =
	// テーブルに追加。
	// ins

	return errors.New("hoge")
}
