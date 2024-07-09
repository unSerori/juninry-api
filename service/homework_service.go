package service

import (
	"errors"
	"fmt"
	"juninry-api/common"
	"juninry-api/dip"
	"juninry-api/model"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

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
// インターフェース型で依存性を受け取ることにより、具体的な実装(gin.Context, GinContextWrapper)ではなくインターフェースに依存し、依存性逆転が実現できる。
func (s *HomeworkService) SubmitHomework(uploader dip.FileUpLoader, bHW model.HomeworkSubmission, form *multipart.Form) error {
	// 画像の保存
	images := form.File["images"] // スライスからimages fieldを取得
	// 保存先ディレクトリ
	dst := "./upload/homework"
	// ディレクトリが存在しない場合
	if _, err := os.Stat(dst); os.IsNotExist(err) { // ファイル情報を取得, 取得できないならerrができる // 取得できなかったとき、ファイルが存在しないことが理由なら新しく作成
		if err := os.MkdirAll(dst, 0644); err != nil {
			return err
		}
	}

	// それぞれのファイルを保存
	for _, image := range images {
		// バリデーション

		// 画像リクエストのContent-Typeから形式(png, jpg, jpeg, gif)の確認
		mimeType := image.Header.Get("Content-Type") // リクエスト画像のmime typeを取得
		ok, _ := validMime(mimeType)                 // 許可されたMIMEタイプか確認
		if !ok {
			return common.NewErr(common.ErrTypeInvalidFileFormat, common.WithMsg("the Content-Type of the request image is invalid"))
		}
		// ファイルのバイナリからMIMEタイプを推測し確認、拡張子を取得
		buffer := make([]byte, 512) // バイトスライスのバッファを作成
		file, err := image.Open()   // multipart.Formを実装するFileオブジェクトを直接取得
		if err != nil {
			return err
		}
		defer file.Close()                                 // 終了後破棄
		file.Read(buffer)                                  // ファイルをバッファに読み込む
		mimeTypeByBinary := http.DetectContentType(buffer) // 読み込んだバッファからコンテントタイプを取得
		ok, validType := validMime(mimeTypeByBinary)       // 許可されたMIMEタイプか確認
		if !ok {
			return common.NewErr(common.ErrTypeInvalidFileFormat, common.WithMsg("the Content-Type inferred from the request image binary is invalid"))
		}
		fileExt := strings.Split(validType, "/")[1] // 画像の種類を取得して拡張子として保存

		// ファイル名をuuidで作成
		fileName, err := uuid.NewRandom() // 新しいuuidの生成
		if err != nil {
			return err
		}

		// 保存
		fmt.Printf("image.Filename: %v\n", image.Filename)     // ファイル名
		fmt.Printf("mimeType: %v\n", mimeType)                 // リクエストヘッダからのContent-Type
		fmt.Printf("mimeTypeByBinary: %v\n", mimeTypeByBinary) // バイナリからのContent-Type
		fmt.Printf("validType: %v\n", validType)
		fmt.Printf("fileExt: %v\n", fileExt)
		fmt.Println("filename: " + dst + "/" + fileName.String() + "." + fileExt)
		//out, err := os.OpenFile()  //uploader.SaveUploadedFile(image, dst+"/"+fileName.String()+"."+fileExt) // c.SaveUploadedFile(image, dst+"/"+fileName.String()+".png")
	}

	// 画像名スライスを文字列に変換し、
	// list :=
	// 画像一覧を提出中間テーブル構造体インスタンスに追加し、
	// bHW.list =
	// テーブルに追加。
	// ins

	return errors.New("hoge")
}

// 許可されたMIMEタイプかどうかを確認、許可されていた場合は一致したタイプを返す
func validMime(mimetype string) (bool, string) {
	// 有効なファイルタイプを定義
	var allowedMimeTypes = []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
		"image/gif",
	}

	for _, allowedMimeType := range allowedMimeTypes {
		if strings.EqualFold(allowedMimeType, mimetype) { // 大文字小文字を無視して文字列比較
			return true, allowedMimeType // 一致した時点で早期リターン
		}
	}

	return false, ""
}
