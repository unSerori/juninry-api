package controller

import (
	"fmt"
	"juninry-api/common/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// json
func TestJson(c *gin.Context) {
	// 成功ログ
	logging.SuccessLog("JSON for testing.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(http.StatusOK, gin.H{ // bodyがJSON形式のレスポンスを返す
		"srvResCode": http.StatusText(resStatusCode),       // メッセージ
		"srvResData": gin.H{"message": "hello go server!"}, // データ
	})
}

func CfmReq(c *gin.Context) {
	fmt.Println("Request confirmed!!!!!!!!!!!!!!!!!!!!")

	fmt.Print("method: ")
	fmt.Println(c.Request.Method)
	fmt.Print("url: ")
	fmt.Println(c.Request.URL)
	// fmt.Print("tls ver: ")
	// fmt.Println(c.Request.TLS.Version)
	fmt.Print("header: ")
	fmt.Println(c.Request.Header)
	fmt.Print("body: ")
	fmt.Println(c.Request.Body)
	fmt.Print("url query: ")
	fmt.Println(c.Request.URL.Query())
	fmt.Println()
}
