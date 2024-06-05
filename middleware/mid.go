package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// ロギング
func MidLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// リクエストを受け取った時のログ
		log.Printf("Received request.\n")                        // リクエストの受理ログ
		log.Printf("Time: %v\n", time.Now())                     // 時刻
		log.Printf("Request method: %s\n", ctx.Request.Method)   // メソッドの種類
		log.Printf("Request path: %s\n\n", ctx.Request.URL.Path) // リクエストパラメータ

		// リクエストを次のハンドラに渡す
		ctx.Next()

		// レスポンスを返した後のログ
		log.Printf("Sent response.\n")                             // レスポンスの送信ログ
		log.Printf("Time: %v\n", time.Now())                       // 時刻
		log.Printf("Response Status: %d\n\n", ctx.Writer.Status()) // ステータスコード
	}
}
