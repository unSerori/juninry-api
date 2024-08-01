package service

import (
	"errors"
	"io"
	"juninry-api/common/logging"
	"juninry-api/model"
	"juninry-api/utility"
	"juninry-api/utility/custom"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type HomeworkService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// 課題の提出履歴の構造体
type SubmissionRecord struct {
	LimitDate      	time.Time	`json:"limitDate"`			// 締め切り
	SubmissionCount int			`json:"submissionCount"`	// 提出数
	HomeworkCount  	int			`json:"homeworkCount"`		// 課題数
}

// 課題の提出履歴を取得
func (s *HomeworkService) GetHomeworkRecord(userId string, targetMonth time.Time) ([]SubmissionRecord, error) {
	// ユーザーが生徒かな
	isJunior, err := model.IsJunior(userId)
	if err != nil {
		return nil, err
	}
	if !isJunior {
		return nil, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// その月締め切りの課題一覧を取得
	// ユーザーIDからクラスを取得
	classMemberships, err := model.FindClassMemberships(userId)
	if err != nil {
		return nil, err
	}

	// クラスIDをスライスに変換
	var classUuids []string
	for _, value := range classMemberships {
		classUuids = append(classUuids, value.ClassUuid)
	}

	// クラスIDから教材一覧を取得
	teachingMaterials, err := model.FindTeachingMaterials(classUuids)
	if err != nil {
		return nil, err
	}

	// 教材IDをスライスに変換
	var materialUuids []string
	for _, value := range teachingMaterials {
		materialUuids = append(materialUuids, value.TeachingMaterialUuid)
	}

	// 教材IDから課題一覧を取得
	homeworks, err := model.FindHomeworks(materialUuids, targetMonth)
	if err != nil {
		return nil, err
	}

	// 課題の日付をキーとしたMap
	var homeworkUuidsMap = make(map[time.Time] []string)

	// 提出期限を1日でまとめたのキーに課題UUIDを追加
	for _, v := range homeworks {
		// 時間を24時間単位に切り捨てる
		homeworkUuidsMap[v.HomeworkLimit.Truncate(24 * time.Hour)] = append(homeworkUuidsMap[v.HomeworkLimit.Truncate(24 * time.Hour)], v.HomeworkUuid)
	}

	// レスポンスの構造体
	var submissionRecord []SubmissionRecord

	for key, value := range homeworkUuidsMap {
		// 課題が提出されているかを確認
		count, err := model.CheckHomeworkSubmission(value)
		if err != nil {
			return nil, err
		}

		// 日付と課題数、提出数をどこどこ追加
		submissionRecord = append(submissionRecord, SubmissionRecord{LimitDate: key, HomeworkCount: len(value), SubmissionCount: int(count)})
	}

	return submissionRecord, nil
}


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

// クラスごとに課題データをまとめた構造体
type ClassHomeworkSummary struct {
	ClassName    string         `json:"className"`    //提出期限
	HomeworkData []HomeworkData `json:"homeworkData"` //課題データのスライス
}

// 宿題登録のリクエストバインド構造体
type BindRegisterHW struct { // model.Homework + classUUID
	model.Homework
	ClassUUID string `json:"classUUID"`
}

// userUuidをuserHomeworkモデルに投げて、受け取ったデータを整形して返す
func (s *HomeworkService) FindHomework(userUuid string) ([]TransformedData, error) {

	// 親には宿題一覧使えないよ
	isPatron, err := model.IsPatron(userUuid)
	if err != nil {
		return nil, err
	}
	if isPatron {	// 親が宿題一覧見ようとしないでね、何も情報とれないんだけどさ、、、
		return nil, custom.NewErr(custom.ErrTypePermissionDenied)
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

// userUuidをuserHomeworkモデルに投げて、次の日が期限の課題データを整形して返す
func (s *HomeworkService) FindClassHomework(userUuid string) ([]ClassHomeworkSummary, error) {

	//user_uuidを絞り込み条件にクソデカ構造体のスライスを受け取る
	userHomeworkList, err := model.FindUserHomeworkforNextday(userUuid)
	if err != nil { //エラーハンドル エラーを上に投げるだけ
		return nil, err
	}

	// クラス名をキー、バリューを課題データのマップにする
	transformedDataMap := make(map[string][]HomeworkData)
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
		transformedDataMap[userHomework.ClassName] = append(transformedDataMap[userHomework.ClassName], homeworkData)
	}

	//作ったマップをさらに整形
	var transformedDataList []ClassHomeworkSummary
	for className, homeworkData := range transformedDataMap {
		transformedData := ClassHomeworkSummary{
			ClassName:    className,
			HomeworkData: homeworkData,
		}
		transformedDataList = append(transformedDataList, transformedData)
	}

	//できたら返す
	return transformedDataList, nil
}

// 宿題登録処理
// インターフェース型で依存性を受け取ることにより、具体的な実装(gin.Context, GinContextWrapper)ではなくインターフェースに依存し、依存性逆転が実現できる。uploader dip.FileUpLoader,
func (s *HomeworkService) SubmitHomework(bHW *model.HomeworkSubmission, form *multipart.Form) error {
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

	// 保存した画像リスト
	var imageNameList []string

	// それぞれのファイルを保存
	for _, image := range images {
		// バリデーション

		// ファイルサイズ
		var maxSize int64                                                              // 上限設定値
		maxSize = 5242880                                                              // default値10MB
		if maxSizeByEnv := os.Getenv("MULTIPART_IMAGE_MAX_SIZE"); maxSizeByEnv != "" { // 空文字でなければ数値に変換する
			var err error
			maxSizeByEnvInt, err := strconv.Atoi(maxSizeByEnv) // 数値に変換
			if err != nil {
				return err
			}
			maxSize = int64(maxSizeByEnvInt) // int64に変換
		}
		if image.Size > maxSize { // ファイルサイズと比較する
			return custom.NewErr(custom.ErrTypeFileSizeTooLarge)
		}

		// 画像リクエストのContent-Typeから形式(png, jpg, jpeg, gif)の確認
		mimeType := image.Header.Get("Content-Type") // リクエスト画像のmime typeを取得
		ok, _ := validMime(mimeType)                 // 許可されたMIMEタイプか確認
		if !ok {
			return custom.NewErr(custom.ErrTypeInvalidFileFormat, custom.WithMsg("the Content-Type of the request image is invalid"))
		}
		// ファイルのバイナリからMIMEタイプを推測し確認、拡張子を取得
		buffer := make([]byte, 512) // バイトスライスのバッファを作成
		file, err := image.Open()   // multipart.Formを実装するFileオブジェクトを直接取得  // このバイナリはファイルタイプの特定とファイル保存書き込み処理で使う
		if err != nil {
			return err
		}
		defer file.Close()                                 // 終了後破棄
		file.Read(buffer)                                  // ファイルをバッファに読み込む  // 読み込んだバイト数とエラーを返す
		mimeTypeByBinary := http.DetectContentType(buffer) // 読み込んだバッファからコンテントタイプを取得
		ok, validType := validMime(mimeTypeByBinary)       // 許可されたMIMEタイプか確認
		if !ok {
			return custom.NewErr(custom.ErrTypeInvalidFileFormat, custom.WithMsg("the Content-Type inferred from the request image binary is invalid"))
		}
		fileExt := strings.Split(validType, "/")[1] // 画像の種類を取得して拡張子として保存

		// ファイル名をuuidで作成
		fileNameWithoutExt, err := uuid.NewRandom() // 新しいuuidの生成
		if err != nil {
			return err
		}
		fileName := fileNameWithoutExt.String() + "." + fileExt // ファイルネームを生成
		filePath := dst + "/" + fileName                        // ファイルパスを生成

		// 保存 //uploader.SaveUploadedFile(image, dst+"/"+fileName.String()+"."+fileExt) // c.SaveUploadedFile(image, dst+"/"+fileName.String()+".png")

		// ファイルを開く
		oFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644) // ファイルが存在しない場合に新規作成|O_CREATEと組み合わせることで同名ファイル存在時にエラーを発生|書き込み専用で開く
		if err != nil {
			return err
		}
		defer oFile.Close() // リソース解放
		// 読み書き位置の設定
		if _, err := file.Seek(0, io.SeekStart); err != nil { // 書き込みたいデータ
			return err
		}
		if _, err := oFile.Seek(0, io.SeekStart); err != nil { // 開いたファイル
			return err
		}
		// データをコピー
		if _, err := io.Copy(oFile, file); err != nil { // io.Copy()はimage<-*multipart.FileHeaderを解釈できないので、バイナリからファイルタイプを特定するために取得したFileオブジェクトを利用
			return nil
		}

		// 保存した画像リストに追加
		imageNameList = append(imageNameList, fileName)
	}

	// 画像名スライスを文字列に変換し、
	imageNameListString := strings.Join(imageNameList, ", ")
	// 画像一覧を提出中間テーブル構造体インスタンスに追加し、
	bHW.ImageNameListString = imageNameListString

	// 提出日時を現在日時に設定
	bHW.SubmissionDate = time.Now()

	// DBに登録
	_, err := model.StoreHomework(bHW)
	if err != nil {
		// XormのORMエラーを仕分ける
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // errをmysqlErrにアサーション出来たらtrue
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反
				return custom.NewErr(custom.ErrTypeUniqueConstraintViolation)
			default: // ORMエラーの仕分けにぬけがある可能性がある
				return custom.NewErr(custom.ErrTypeOtherErrorsInTheORM)
			}
		}
		// 通常の処理エラー
		return err
	}

	return nil
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

// 宿題登録
func (s *HomeworkService) RegisterHWService(bHW BindRegisterHW, userId string) (string, error) {
	// ユーザー権限の確認
	isTeacher, err := model.IsTeacher(userId)
	if err != nil {
		return "", err
	}
	if !isTeacher { // 教師権限を持っていないならエラー
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return "", custom.NewErr(custom.ErrTypePermissionDenied)
	}
	logging.SuccessLog("User creation authority confirmation complete")

	// 指定されたクラスIDに投稿ユーザー自身が所属しているかを確認
	isMember, err := model.CheckUserClassMembership(bHW.ClassUUID, userId)
	if err != nil {
		return "", err
	}
	if !isMember {
		return "", custom.NewErr(custom.ErrTypePermissionDenied)
	}
	// 投稿者ID追加
	bHW.HomeworkPosterUuid = userId
	logging.SuccessLog("Confirmation of user's affiliation authority complete.")

	// 一意ID生成
	newId, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {
		return "", err
	}
	bHW.HomeworkUuid = newId.String() // 設定

	// 構造体をテーブルモデルに変換
	var hw model.Homework // 構造体のインスタンス
	utility.ConvertStructCopyMatchingFields(&bHW, &hw)

	// 登録
	err = model.CreateHW(hw)
	if err != nil {
		logging.ErrorLog("Failed to register homework", err)
		return "", err
	}

	return bHW.HomeworkUuid, nil
}
