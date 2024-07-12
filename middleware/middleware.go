package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"juninry-api/auth"
	"juninry-api/logging"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ロギング
func LoggingMid() gin.HandlerFunc {
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

// トークン検証
func MidAuthToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ヘッダーからトークンを取得
		headerAuthorization := ctx.GetHeader("Authorization")
		if headerAuthorization == "" { // ヘッダーが存在しない場合
			// エラーログ
			logging.ErrorLog("Authentication unsuccessful.", nil)
			// レスポンス
			ctx.JSON(http.StatusBadRequest, gin.H{
				"srvResCode": 7001,                           // コード
				"srvResMsg":  "Authentication unsuccessful.", // メッセージ
				"srvResData": gin.H{},                        // データ
			})
			ctx.Abort() // 次のルーティングに進まないよう処理を止める。
			return      // 早期リターンで終了
		}

		// トークンの解析を行う。
		token, id, err := auth.ParseToken(headerAuthorization)
		if err != nil {
			// エラーログ
			logging.ErrorLog("Authentication unsuccessful. Maybe that user does not exist. Failed to parse token.", err)
			// レスポンス
			ctx.JSON(http.StatusBadRequest, gin.H{
				"srvResCode": 7008,                                                                                  // コード
				"srvResMsg":  "Authentication unsuccessful. Maybe that user does not exist. Failed to parse token.", // メッセージ
				"srvResData": gin.H{},                                                                               // データ
			})
			ctx.Abort() // 次のルーティングに進まないよう処理を止める。
			return      // 早期リターンで終了
		}

		ctx.Set("token", token) // トークンをコンテキストにセットする。  // _ = token // トークンを破棄。
		ctx.Set("id", id)       // 送信元クライアントのtokenのidを保持

		ctx.Next() // エンドポイントの処理に移行
	}
}

// 同時に一人しか実行させないよ〜
func SingleExecutionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var lock sync.Mutex
		lock.Lock()
		defer lock.Unlock()

		c.Next() // エンドポイントの処理に移行
	}
}

// リクエストボディ容量を制限するミドルウェア
func LimitReqBodySize(maxBytesSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Printf("maxBytesSize: %v\n", maxBytesSize)
		// Content-Lengthからの確認
		if ctx.Request.ContentLength > maxBytesSize {
			// エラーログ
			logging.ErrorLog("Multipart request bodies are too big.", errors.New("request size "+strconv.Itoa(int(maxBytesSize))+"bytes"))
			// レスポンス
			resStatusCode := http.StatusRequestEntityTooLarge
			ctx.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
			ctx.Abort() // リクエスト処理を中止
			return
		}

		// リクエストのボディを実際に読み込んで判定
		buf := make([]byte, maxBytesSize)             // 制限の分だけ読み込めるバッファを用意
		n, err := io.ReadFull(ctx.Request.Body, buf)  // バッファの容量分だけ読み込む  // nは読み込めたバイト数
		if err != nil && err != io.ErrUnexpectedEOF { // EOF以外のエラーが発生した場合は内部エラー
			// エラーログ
			logging.ErrorLog("Internal Server Error.", err)
			// レスポンス
			resStatusCode := http.StatusInternalServerError
			ctx.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
			ctx.Abort() // リクエスト処理を中止
			return
		}
		if int64(n) == maxBytesSize && err == nil { // 制限サイズと同じまで読み込めてしまったら413
			// エラーログ
			logging.ErrorLog("Payload Too Large.", err)
			// レスポンス
			resStatusCode := http.StatusRequestEntityTooLarge
			ctx.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
			ctx.Abort() // リクエスト処理を中止
			return
		}

		// 読み取ったデータをリクエストボディに再設定
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(buf[:n]))

		ctx.Next() // エンドポイントの処理に移行
	}
}
