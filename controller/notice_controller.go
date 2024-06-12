package controller

import (
	"github.com/gin-gonic/gin"
	"juninry-api/logging"
	"juninry-api/service"
	"net/http"
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
