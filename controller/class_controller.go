package controller

import (
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var ClassService = service.ClassService{}

func RegisterClassHandler(c *gin.Context) {
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
	// idAdjusted := id.(string) // アサーション

	print(id)
	time.Sleep(2 * time.Second)

	//構造体に値をバインド
	var bClass model.Class
	if err := c.ShouldBindJSON(&bClass); err != nil {
		// エラーログ
	}
	return

	// 登録処理を投げてなんかいろいろもらう
	// class, err := ClassService.PermissionCheckedClassCreation(idAdjusted, bClass)
	// if err != nil {
	// 	return
	// }

	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": "check",
	})

}

// サービスに登録処理を投げる
