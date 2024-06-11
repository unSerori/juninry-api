package controller

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResCode": 7013,                    // コード
			"srvResMsg":  "The id is not stored.", // メッセージ
			"srvResData": gin.H{},                 // データ
		})
		return
	}

	//問い合わせ処理と失敗レスポンス
	homeworkList, err := homeworkService.FindHomework(id.(string))
	if err != nil { //エラーハンドル
		// エラーログ
		logging.ErrorLog("SQL query failed.", err)
		//レスポンス
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get homework list.")
	// レスポンス
	c.JSON(http.StatusOK, gin.H{
		"srvResData": homeworkList,
	})
}
