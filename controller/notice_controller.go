package controller

import (
	"juninry-api/logging"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var noticeService = service.NoticeService{} // サービスの実体を作る。

// お知らせ1件取得
func GetNoticeDetailHandler(ctx *gin.Context) {

	//notice_uuidの取得
	noticeUuid := ctx.Param("notice_uuid")

	//お知らせのレコードを取ってくる
	noticeDetail, err := noticeService.GetNoticeDetail(noticeUuid)
	if err != nil {
		// エラーログ
		logging.ErrorLog("notice find error", err)
		// レスポンス
		ctx.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	// 成功ログ
	logging.SuccessLog("Successful noticeDetail get.")
	// レスポンス(StatusOK　成功200番)
	ctx.JSON(http.StatusOK, gin.H{
		"srvResMsg":  "Successful noticeDetail get.",
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
