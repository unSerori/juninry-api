package controller

import (
	"errors"
	"fmt"
	"juninry-api/common/logging"
	"juninry-api/service"
	"juninry-api/utility/custom"
	"net/http"

	"github.com/gin-gonic/gin"
)

var nyariotSarvice = service.NyariotSarvice{} // サービスの実体を作る。

// 　その日初めてのログインの時スタンプを付与する
func LoginStampHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	stamp, err := nyariotSarvice.AddLoginStamp(idAdjusted)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypePermissionDenied: // 生徒じゃないので閲覧権限無し
				// エラーログ
				logging.ErrorLog("Forbidden.", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		} else { // カスタムエラー以外の処理エラー
			// エラーログ
			logging.ErrorLog("Internal Server Error.", err)
			// レスポンス
			resStatusCode := http.StatusInternalServerError
			c.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
		}
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("The read process was completed.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": stamp,
	})
}

// 現在のスタンプ数を取得
func GetStampsHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// TODO:返ってくる値に不必要な値があるけど、時間がないのでそのままです。ごめんなさい。
	stamp, err := nyariotSarvice.GetStamp(idAdjusted)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypePermissionDenied: // 生徒じゃないので閲覧権限無し
				// エラーログ
				logging.ErrorLog("Forbidden.", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		} else { // カスタムエラー以外の処理エラー
			// エラーログ
			logging.ErrorLog("Internal Server Error.", err)
			// レスポンス
			resStatusCode := http.StatusInternalServerError
			c.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
		}
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("The read process was completed.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": stamp,
	})
}

// func GetGachaByStampHandler(c *gin.Context) {
// 	// ユーザ特定
// 	id, _ := c.Get("id")
// 	idAdjusted := id.(string) // アサーション

// 	gacha, err := nyariotSarvice.GetStampGacha(idAdjusted)
// 	if err != nil { // エラーハンドル
// 		// カスタムエラーを仕分ける
// 		var customErr *custom.CustomErr
// 		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
// 			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
// 			case custom.ErrTypePermissionDenied: // 生徒じゃないので閲覧権限無し
// 				// エラーログ
// 				logging.ErrorLog("Forbidden.", err)
// 				// レスポンス
// 				resStatusCode := http.StatusForbidden
// 				c.JSON(resStatusCode, gin.H{
// 					"srvResMsg":  http.StatusText(resStatusCode),
// 					"srvResData": gin.H{},
// 				})
// 			default: // カスタムエラーの仕分けにぬけがある可能性がある
// 				// エラーログ
// 				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
// 				// レスポンス
// 				resStatusCode := http.StatusBadRequest
// 				c.JSON(resStatusCode, gin.H{
// 					"srvResMsg":  http.StatusText(resStatusCode),
// 					"srvResData": gin.H{},
// 				})
// 			}
// 		} else { // カスタムエラー以外の処理エラー
// 			// エラーログ
// 			logging.ErrorLog("Internal Server Error.", err)
// 			// レスポンス
// 			resStatusCode := http.StatusInternalServerError
// 			c.JSON(resStatusCode, gin.H{
// 				"srvResMsg":  http.StatusText(resStatusCode),
// 				"srvResData": gin.H{},
// 			})
// 		}
// 		return
// 	}

// 	// 処理後の成功
// 	// 成功ログ
// 	logging.SuccessLog("The read process was completed.")
// 	// レスポンス
// 	resStatusCode := http.StatusOK
// 	c.JSON(resStatusCode, gin.H{
// 		"srvResMsg":  http.StatusText(resStatusCode),
// 		"srvResData": gacha,
// 	})

// }

// 所持アイテム一覧取得(図鑑表示のため持っていないアイテムも帰ってくる)
func GetUserItemBoxHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	items, err := nyariotSarvice.GetUserItems(idAdjusted)
	// 取得できなかった時のエラーを判断
	if err != nil {
		// 処理で発生したエラーのうちカスタムエラーのみ
		var serviceErr *custom.CustomErr
		if errors.As(err, &serviceErr) {
			switch serviceErr.Type {
			case custom.ErrTypePermissionDenied:
				// エラーログ(権限無し)
				logging.ErrorLog("Do not have the necessary permissions", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			default:
				// エラーログ(権限無し)
				logging.ErrorLog("StatusBadRequest", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		}
		// エラーログ
		logging.ErrorLog("itme find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful items get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful items get.",
		"srvResData": items,
	})
}

// アイテム詳細取得
func GetItemDetail(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// アイテムUUIDを取得
	itemUuid := c.Param("item_uuid")

	item, err := nyariotSarvice.GetItemDetail(idAdjusted, itemUuid)
	// 取得できなかった時のエラーを判断
	if err != nil {
		// 処理で発生したエラーのうちカスタムエラーのみ
		var serviceErr *custom.CustomErr
		if errors.As(err, &serviceErr) {
			switch serviceErr.Type {
			case custom.ErrTypePermissionDenied:
				// エラーログ(権限無し)
				logging.ErrorLog("Do not have the necessary permissions", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			default:
				// エラーログ(権限無し)
				logging.ErrorLog("StatusBadRequest", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		}
		// エラーログ
		logging.ErrorLog("itme find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful items get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful items get.",
		"srvResData": item,
	})

}
