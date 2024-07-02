package controller

import (
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
		// エラーログ
		logging.ErrorLog("Failed to retrieve image request.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 提出記録処理と失敗レスポンス
	err = homeworkService.SubmitHomework(c, bHW, form)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		return
	}

	// 成功レスポンス 200番
	// 成功ログ
	logging.SuccessLog("Successful submission homework.")
	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{},
	})
}
