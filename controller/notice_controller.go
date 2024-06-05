package controller

import (
	"juninry-api/logging"
	"juninry-api/model"
	//"log"
	"github.com/gin-gonic/gin"
)

func Get_Notices_Handler(ctx *gin.Context) {
	// 結果を格納する変数(findの結果)
	get_notices := []model.Notice{}
	//テストコード(全件取得)
	err := model.DBInstance().Find(
			&get_notices,
	)

	// 取得できなかった時のエラーを判断
	if err != nil {
		// エラーログ
		logging.ErrorLog("notice find error", err)
		// レスポンス(500番台はサーバー側の失敗)
		ctx.JSON(500, gin.H{
			"srvResCode": 0000,
			"srvResMsg":  "notice find error.",
			"srvResData": gin.H{},
		})
		return    //　<-返すよって型指定してないから切り上げるだけ
	}

	//確認用
	//log.Println(get_notices)

	// レスポンス(200番台はサーバー側の失敗)
	ctx.JSON(200, gin.H{
		"srvResCode": 1130,
		"srvResMsg":  "Successful get notice.",
		"srvResData": gin.H{
			"notices" : get_notices,
		},
	})

}


// 関数名の先頭が大文字の場合、pubulic