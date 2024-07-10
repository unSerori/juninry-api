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

var OuchiService = service.OuchiService{}

func RegisterOuchiHandler(ctx *gin.Context) {
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

	//構造体に値をバインド
	var bOuchi model.Ouchi
	if err := ctx.ShouldBindJSON(&bOuchi); err != nil {
		fmt.Print("バインド失敗")
		// エラーログ
		return
	}

	// 登録処理を投げてなんかいろいろもらう
	ouchi, err := OuchiService.PermissionCheckedOuchiCreation(idAdjusted, bOuchi)
	if err != nil {
		var serviceErr *common.CustomErr
		if errors.As(err, &serviceErr) { // カスタムエラーの場合
			switch serviceErr.Type {
			case common.ErrTypePermissionDenied: // 権限を持っていない
				logging.ErrorLog("Do not have the necessary permissions", err)
				resStatusCode := http.StatusForbidden
				ctx.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			case common.ErrTypeMaxAttemptsReached: // 最大試行数を超えた
				logging.ErrorLog("Max attempts reached", err)
			}
		} else {
			// エラーログ
			logging.ErrorLog("Ouchi creation was not possible due to other problems.", err)
		}
		resStatusCode := http.StatusBadRequest
		ctx.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// レスポンス
	resStatusCode := http.StatusCreated
	ctx.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": ouchi,
	})
}
