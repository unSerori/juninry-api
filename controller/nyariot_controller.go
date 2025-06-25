package controller

import (
	"errors"
	"fmt"
	"juninry-api/common/custom"
	"juninry-api/common/logging"
	"juninry-api/service"
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

// スタンプでガチャ
func GetGachaByStampHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	gacha, err := nyariotSarvice.GetStampGacha(idAdjusted)
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
				return
			case custom.ErrTypeResourceUnavailable: // アイテムないよエラー
				// エラーログ
				logging.ErrorLog("don't own the item.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{ // アイテムないよエラー
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			case custom.ErrTypeUnforeseenCircumstances:	// ガチャ回すのに必要なスタンプないよ
				//エラーログ
				logging.ErrorLog("unforeseen circumstances", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return


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
		"srvResData": gacha,
	})

}

// ポイントでガチャ
func GetGachaByPointHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// まわす回数を取得
	count := c.Param("count")

	gacha, err := nyariotSarvice.GetPointGacha(idAdjusted, count)
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
				return
			case custom.ErrTypeResourceUnavailable: // アイテムないよエラー
				// エラーログ
				logging.ErrorLog("don't own the item.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{ // アイテムないよエラー
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			case custom.ErrTypeUnforeseenCircumstances:
				//エラーログ
				logging.ErrorLog("unforeseen circumstances", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return

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
		"srvResData": gacha,
	})

}

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

// 所持ニャリオット一覧取得(図鑑表示のため持っていないアイテムも帰ってくる)
func GetUserNyariotInventoryHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	nyairots, err := nyariotSarvice.GetUserNyariots(idAdjusted)
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
		logging.ErrorLog("nyariots find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful nyariots get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful nyariots get.",
		"srvResData": nyairots,
	})
}

// ニャリオット詳細取得
func GetNyariotDetail(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// ニャリオットUUIDを取得
	nyariotUuid := c.Param("nyariot_uuid")

	nyariot, err := nyariotSarvice.GetNyariotDetail(idAdjusted, nyariotUuid)
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
		logging.ErrorLog("nyariot find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful nyariot get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful nyariot get.",
		"srvResData": nyariot,
	})

}

// ニャリオット更新
func ChangeMainNariot(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// ニャリオットUUIDを取得
	nyariotUuid := c.Param("nyariot_uuid")

	err := nyariotSarvice.ChangeNariot(idAdjusted, nyariotUuid)
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
		logging.ErrorLog("nyariot find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Main Nyariot Changed.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Main Nyariot Changed.",
		"srvResData": gin.H{},
	})

}

// メインニャリオット取得
func GetMainNyariotHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	nyariot, err := nyariotSarvice.GetMainNyariot(idAdjusted)
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
		logging.ErrorLog("nyariot find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful nyariot get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful nyariot get.",
		"srvResData": nyariot,
	})
}

// 空腹度取得
func GetHungryStatusHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	nyariot, err := nyariotSarvice.GetHungryStatus(idAdjusted)
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
		logging.ErrorLog("nyariot find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful nyariot SatityDegrees get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful nyariot SatityDegrees get.",
		"srvResData": nyariot,
	})
}

// 空腹度更新
func UpdateHungryStatusHandler(c *gin.Context) {
	// ユーザ特定
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// アイテムUUIDを取得
	itemUuid := c.Param("item_uuid")

	hungryStatus, err := nyariotSarvice.UpdateHungryStatus(idAdjusted, itemUuid)
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
			case custom.ErrTypeResourceUnavailable: // アイテムないよエラー
				// エラーログ
				logging.ErrorLog("don't own the item.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{ // アイテムないよエラー
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
		logging.ErrorLog("nyariot find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		c.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// 成功ログ
	logging.SuccessLog("Successful nyariot SatityDegrees get.")
	// レスポンス(StatusOK　成功200番)
	c.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful nyariot SatityDegrees get.",
		"srvResData": hungryStatus ,
	})

}
