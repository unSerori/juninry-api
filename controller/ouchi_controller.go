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
	"github.com/go-sql-driver/mysql"
)

var OuchiService = service.OuchiService{}

// おうち新規作成
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

// おうちの招待コード更新
func GenerateOuchiInviteCodeHandler(ctx *gin.Context) {
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

	// おうちUUIDを取得
	ouchiUuid := ctx.Param("ouchi_uuid")

	// 招待コード登録します
	ouchi, err := OuchiService.PermissionCheckedRefreshOuchiInviteCode(idAdjusted, ouchiUuid)
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
			case common.ErrTypeNoResourceExist: // リソースがない
				logging.ErrorLog("The resource does not exist", err)
				resStatusCode := http.StatusNotFound
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

// おうち参加処理
func JoinOuchiHandler(c *gin.Context) {
	// ユーザーを特定する
	id, exists := c.Get("id")
	if !exists { // idがcに保存されていない。
		// エラーログ
		logging.ErrorLog("The id is not stored.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	idAdjusted := id.(string) // アサーション
	// 招待コードを取得
	inviteCode := c.Param("invite_code")

	// おうちに参加
	ouchiName, err := OuchiService.PermissionCheckedJoinOuchi(idAdjusted, inviteCode)
	if err != nil {
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // 第一引数のerrが第二引数の型にキャスト可能ならキャストしてtrue
			if mysqlErr.Number == 1062 { // 重複エラー
				logging.ErrorLog("The ouchi has already joined", err)
				resStatusCode := http.StatusConflict
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			}
		}

		// いろいろ長いからややこしいけど、swich文でエラー振り分けてるだけ
		var serviceErr *common.CustomErr
		if errors.As(err, &serviceErr) { // カスタムエラーの場合
			switch serviceErr.Type {
			case common.ErrTypePermissionDenied: // 権限を持っていない
				logging.ErrorLog("Do not have the necessary permissions", err)
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			case common.ErrTypeNoResourceExist: // 招待コード違います
				logging.ErrorLog("The resource does not exist", err)
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default:
				logging.ErrorLog("Ouchi creation was not possible due to other problems.", err)
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			}

		}
	}

	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"ouchiName": ouchiName,
		},
	})

}
