// ビジネスロジック

package domain

import (
	"juninry-api/common"
	"juninry-api/logging"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ビジネスロジックの構造体  // infrastructure層の処理を呼び出したいが、このままではdomain<->infrastructure間で循環参照してしまう。そのため実際の実装(infra)に依存するのではなく、提供元の同レイヤー内のrepository interface(実装された処理関数の型と呼べる)を利用することで、repositoryを介して具体的な実装を利用できる。
type TeachingMaterialLogic struct {
	r TeachingMaterialRepository
}

// ファクトリー関数
func NewTeachingMaterialLogic(r TeachingMaterialRepository) *TeachingMaterialLogic {
	return &TeachingMaterialLogic{r: r}
}

// ビジネスロジック  // アプリケーション層からのエンティティへの影響をビジネスルールに従って実現する

// 教科作成用の構造体を作るために、ユーザーidとクラスidの整合性をとる
func (l *TeachingMaterialLogic) IntegrityStruct(id string, bTM TeachingMaterial) (TeachingMaterial, error) {
	// ユーザーの権限の確認
	permission, isFound, err := l.r.GetPermissionInfoById(id)
	if err != nil || !isFound {
		return TeachingMaterial{}, err
	}
	if permission != 1 { // teacher, junior, patron
		return TeachingMaterial{}, common.NewErr(common.ErrTypePermissionDenied)
	}

	// ユーザーがクラスに属していることを確認
	isFound, err = l.r.IsUserInClass(bTM.ClassUuid, id)
	if err != nil || !isFound {
		return TeachingMaterial{}, err
	}
	if permission != 1 { // teacher, junior, patron
		return TeachingMaterial{}, common.NewErr(common.ErrTypePermissionDenied) // teacherであっても対象クラスに属していなければ権限なし。
	}

	// 教科があるか
	isFound, err = l.r.IsSubjectExists(bTM.SubjectId)
	if err != nil || !isFound {
		return TeachingMaterial{}, err
	}
	if permission != 1 { // teacher, junior, patron
		return TeachingMaterial{}, common.NewErr(common.ErrTypeNoResourceExist)
	}

	// 問題ないので構造体を返す
	return bTM, nil
}

// 画像保存
// 返り血はファイル名(:拡張子を含まない)
func (l *TeachingMaterialLogic) ValidateImage(form *multipart.Form) (string, error) {
	// filesスライスからimage fieldsのひとつめを取得
	image := form.File["image"][0]

	// 保存先ディレクトリの確保
	dst := "./upload/homework"
	l.r.CreateDstDir("./upload/t_material", 0644)

	// バリデーション

	// ファイルサイズの制限(context由来)
	var maxSize int64                                                              // 上限設定値
	maxSize = 5242880                                                              // default値10MB
	if maxSizeByEnv := os.Getenv("MULTIPART_IMAGE_MAX_SIZE"); maxSizeByEnv != "" { // 空文字でなければ数値に変換する
		var err error
		maxSizeByEnvInt, err := strconv.Atoi(maxSizeByEnv) // 数値に変換
		if err != nil {
			return "", err
		}
		maxSize = int64(maxSizeByEnvInt) // int64に変換
	}
	if image.Size > maxSize { // ファイルサイズと比較する
		return "", common.NewErr(common.ErrTypeFileSizeTooLarge)
	}
	// ファイルサイズの制限(Content-Length, binary由来)
	// 画像リクエストのContent-Typeから形式(png, jpg, jpeg, gif)の確認
	mimeType := image.Header.Get("Content-Type") // リクエスト画像のmime typeを取得
	ok, _ := validMime(mimeType)                 // 許可されたMIMEタイプか確認
	if !ok {
		return "", common.NewErr(common.ErrTypeInvalidFileFormat, common.WithMsg("the Content-Type of the request image is invalid"))
	}
	// ファイルのバイナリからMIMEタイプを推測し確認、拡張子を取得
	buffer := make([]byte, 512) // バイトスライスのバッファを作成
	file, err := image.Open()   // multipart.Formを実装するFileオブジェクトを直接取得  // このバイナリはファイルタイプの特定とファイル保存書き込み処理で使う
	if err != nil {
		return "", err
	}
	defer file.Close()                                 // 終了後破棄
	file.Read(buffer)                                  // ファイルをバッファに読み込む  // 読み込んだバイト数とエラーを返す
	mimeTypeByBinary := http.DetectContentType(buffer) // 読み込んだバッファからコンテントタイプを取得
	ok, validType := validMime(mimeTypeByBinary)       // 許可されたMIMEタイプか確認
	if !ok {
		return "", common.NewErr(common.ErrTypeInvalidFileFormat, common.WithMsg("the Content-Type inferred from the request image binary is invalid"))
	}
	fileExt := strings.Split(validType, "/")[1] // 画像の種類を取得して拡張子として保存

	// ファイル名をuuidで作成
	fileNameWithoutExt, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {
		return "", err
	}
	fileName := fileNameWithoutExt.String() + "." + fileExt // ファイルネームを生成
	filePath := dst + "/" + fileName                        // ファイルパスを生成

	// 保存
	err = l.r.UpLoadImage(filePath, file)
	if err != nil {
		return "", err
	}

	return fileNameWithoutExt.String(), nil
}

// 教材構造体をインサート
func (l *TeachingMaterialLogic) CreateTM(tm TeachingMaterial) error {
	err := l.r.CreateTM(tm)
	if err != nil {
		return err
	}

	return nil
}

// 分離された処理関数

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
			logging.InfoLog("True validMime", "True validMime/mimetype: "+mimetype+", allowedMimeType: "+allowedMimeType)
			return true, allowedMimeType // 一致した時点で早期リターン
		}
	}

	logging.InfoLog("False validMime", "False validMime/mimetype: "+mimetype)

	return false, ""
}
