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

	//TODO: 取得する名前わかりません
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

//ユーザの所属するクラスのお知らせ全件取得
func GetAllNoticesHandler(ctx *gin.Context) {
	// 絞り込み条件
	userUuid := "3cac1684-c1e0-47ae-92fd-6d7959759224"

	// userUuidからお知らせ一覧を持って来る(厳密にはserviceにuserUuidを渡す)
	notices, err := noticeService.FindAllNotices(userUuid)
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

