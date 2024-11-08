// レスポンスの一連処理

package responder

import (
	"fmt"
	"juninry-api/common/custom"
	"juninry-api/common/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 標準的かつ柔軟（:設定値を全て呼び出し側で設定する必要がある）なJSONレスポンス
// msg: 先頭大文字コロンあり
// err: err or 先頭小文字コロンなし or nil
func SendJSON(ctx *gin.Context, code int, msg string, err error, body gin.H) {
	// 失敗/成功ログ
	if err != nil {
		logging.ErrorLog(msg, err)
	} else {
		logging.SuccessLog(msg)
	}
	// レスポンス
	resStatusCode := code
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": body,
	})
}

// マッピング失敗
// err: err or 先頭小文字コロンなし
func SendFailedBindJSON(ctx *gin.Context, err error) {
	// エラーログ
	logging.ErrorLog("Failure to bind request.", err)
	// レスポンス
	resStatusCode := http.StatusBadRequest
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{},
	})
}

// サービス処理の失敗レスポンス
// err: err or 先頭小文字コロンなし
func SendFailedService(ctx *gin.Context, code int, err error) {
	// エラーログ
	logging.ErrorLog(http.StatusText(code)+".", err)
	// レスポンス
	resStatusCode := code
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{},
	})
}

// サービス処理のカスタムエラー仕分けのdefault
func SendFailedServiceDefault(ctx *gin.Context, code int, customErr *custom.CustomErr, err error) {
	// エラーログ
	logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
	// レスポンス
	resStatusCode := code
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{},
	})
}

// 成功レスポンス
func SendSuccess(ctx *gin.Context, code int, body interface{}) {
	// 成功ログ
	logging.SuccessLog(http.StatusText(code) + ".")
	// レスポンス
	resStatusCode := code
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": body,
	})
}
