package controller

import (
	"fmt"
	"juninry-api/logging"
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
	// form fields
	hwId := c.PostForm("homeworkUUID")
	fmt.Printf("hwId: %v\n", hwId)
	logging.SimpleLog("log test\n")
}
