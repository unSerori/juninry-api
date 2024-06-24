package controller

import (
	"fmt"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var homeworkService = service.HomeworkService{}

// 課題全件取得
func FindHomeworkHandler(c *gin.Context) {
	// ユーザーを特定する
	id, exists := c.Get("id")
	if !exists { // idがcに保存されていない。
		// エラーログ
		logging.ErrorLog("The id is not stored.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	idAdjusted := id.(string) // アサーション

	//問い合わせ処理と失敗レスポンス
	homeworkList, err := homeworkService.FindHomework(idAdjusted)
	if err != nil { //エラーハンドル
		// エラーログ
		logging.ErrorLog("SQL query failed.", err)
		//レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get homework list.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": homeworkList,
	})
}

// 宿題提出
func SubmitHomeworkHandler(c *gin.Context) {
	// form fields 構造体にマッピング
	var bHW model.HomeworkSubmission     // 構造体のインスタンス
	if err := c.Bind(&bHW); err != nil { // フォームフィールドの直接取得  hwId := c.PostForm("homeworkUUID")
		// エラーログ
		logging.ErrorLog("Failure to bind request.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// form files取得
	form, err := c.MultipartForm() // フォームを取得
	if err != nil {
		return
	}
	images := form.File["images"] // スライス
	// 保存先ディレクトリ
	dst := "./upload/homework"
	// それぞれのファイルを保存
	for _, image := range images {
		fmt.Printf("image.Filename: %v\n", image.Filename)
		// ファイル名をuuidで作成
		fileName, err := uuid.NewRandom() // 新しいuuidの生成
		if err != nil {
			return
		}
		// バリデーション
		// TODO: 形式(png, jpg, jpeg, gif, HEIF)
		// TODO: ファイルの種類->拡張子
		// TODO: パーミッション
		// 保存
		c.SaveUploadedFile(image, dst+"/"+fileName.String()+".png")
	}
	c.JSON(200, gin.H{})
}
