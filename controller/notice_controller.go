package controller

import (
	"errors"
	"fmt"
	"juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var noticeService = service.NoticeService{} // サービスの実体を作る。

// 新規お知らせ登録
func RegisterNoticeHandler(ctx *gin.Context) {

	// 構造体にマッピング
	var bNotice model.Notice // 構造体のインスタンス
	if err := ctx.ShouldBindJSON(&bNotice); err != nil {
		// エラーログ
		logging.ErrorLog("Failure to bind request.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		ctx.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// ユーザーを特定する(ctxに保存されているidを取ってくる)
	id, exists := ctx.Get("id")
	if !exists { // idがcに保存されていない。
		// エラーログ
		logging.ErrorLog("The id is not stored.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		ctx.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	idAdjusted := id.(string) // アサーション
	bNotice.UserUuid = idAdjusted

	// 登録処理と失敗レスポンス
	err := noticeService.RegisterNotice(bNotice)
	if err != nil { // エラーハンドル
		// エラータイプを定義
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypePermissionDenied: // 非管理者ユーザーの場合
				// エラーログ
				logging.ErrorLog("Forbidden.", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
				// レスポンス
				resStatusCode := http.StatusBadRequest
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		} else { // カスタムエラー以外の処理エラー
			// エラーログ
			logging.ErrorLog("Internal Server Error.", err)
			// レスポンス
			resStatusCode := http.StatusInternalServerError
			ctx.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
		}
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful notice registration.")
	// レスポンス
	resStatusCode := http.StatusOK
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{},
	})

}

// お知らせ1件取得
func GetNoticeDetailHandler(ctx *gin.Context) {

	//notice_uuidの取得
	noticeUuid := ctx.Param("notice_uuid")

	//お知らせのレコードを取ってくる
	noticeDetail, err := noticeService.GetNoticeDetail(noticeUuid)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeNoResourceExist: // リソースがなく見つからない
				// エラーログ
				logging.ErrorLog("Not Found.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
				// レスポンス
				resStatusCode := http.StatusBadRequest
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		} else { // カスタムエラー以外の処理エラー
			// エラーログ
			logging.ErrorLog("Internal Server Error.", err)
			// レスポンス
			resStatusCode := http.StatusInternalServerError
			ctx.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
		}
		return
	}

	// 成功ログ
	logging.SuccessLog("Successful noticeDetail get.")
	// レスポンス(StatusOK　成功200番)
	resStatusCode := http.StatusOK
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": noticeDetail,
	})
}

// ユーザの所属するクラスのお知らせ全件取得
func GetAllNoticesHandler(ctx *gin.Context) {
	// 絞り込み条件
	// ユーザーを特定する
	id, exists := ctx.Get("id")
	if !exists { // idがcに保存されていない。
		// エラーログ
		logging.ErrorLog("The id is not stored.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		ctx.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	idAdjusted := id.(string) // アサーション
	

	// userUuidからお知らせ一覧を持って来る(厳密にはserviceにuserUuidを渡す)
	notices, err := noticeService.FindAllNotices(idAdjusted)
	// 取得できなかった時のエラーを判断
	if err != nil {
		// 処理で発生したエラーのうちカスタムエラーのみ
		var serviceErr *common.CustomErr
		if errors.As(err, &serviceErr) {
				switch serviceErr.Type {
				case common.ErrTypePermissionDenied :
						// エラーログ(権限無し)
						logging.ErrorLog("Do not have the necessary permissions", err)
						// レスポンス
						resStatusCode := http.StatusForbidden
						ctx.JSON(resStatusCode, gin.H{
							"srvResMsg":  http.StatusText(resStatusCode),
							"srvResData": gin.H{},
						})
						return
					default: 
					// エラーログ(権限無し)
					logging.ErrorLog("aiueos", err)
					// レスポンス
					resStatusCode := http.StatusBadRequest
					ctx.JSON(resStatusCode, gin.H{
						"srvResMsg":  http.StatusText(resStatusCode),
						"srvResData": gin.H{},
					})
				}
			}
		// エラーログ
		logging.ErrorLog("notice find error", err)
		// レスポンス(StatusInternalServerError サーバーエラー500番)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"srvResData": gin.H{},
		})
		return //　<-返すよって型指定してないから切り上げるだけ
	}

	// レスポンス(StatusOK　成功200番)
	ctx.JSON(http.StatusOK, gin.H{
		"srvResData": gin.H{
			"notices": notices,
		},
	})
}

// お知らせ既読処理
func NoticeReadHandler(ctx *gin.Context) {

	// ユーザーを特定する
	id, exists := ctx.Get("id")
	if !exists { // idがcに保存されていない。
		// エラーログ
		logging.ErrorLog("The id is not stored.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		ctx.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	idAdjusted := id.(string) // アサーション

	//notice_uuidの取得
	noticeUuid := ctx.Param("notice_uuid")

	// 構造体にマッピング
	bRead := model.NoticeReadStatus{
		NoticeUuid: noticeUuid,
		UserUuid:   idAdjusted,
	}

	// 登録処理と失敗レスポンス
	err := noticeService.ReadNotice(bRead)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeUniqueConstraintViolation: // 一意性制約違反
				// エラーログ
				logging.ErrorLog("Conflict.", err)
				// レスポンス
				resStatusCode := http.StatusConflict
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case common.ErrTypePermissionDenied: // 権限なし
				logging.ErrorLog("Do not have the necessary permissions", err)
				resStatusCode := http.StatusForbidden
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
				// レスポンス
				resStatusCode := http.StatusBadRequest
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		} else { // カスタムエラー以外の処理エラー
			// エラーログ
			logging.ErrorLog("Internal Server Error.", err)
			// レスポンス
			resStatusCode := http.StatusInternalServerError
			ctx.JSON(resStatusCode, gin.H{
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
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{
			//TODO:返すものがあるなら入れる
		},
	})

}
