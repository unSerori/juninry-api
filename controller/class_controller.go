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

var ClassService = service.ClassService{}

// クラス一覧取得
func GetAllClassesHandler(c *gin.Context) {
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

	// サービスに投げるよ
	classes, err := ClassService.GetClassList(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get class list.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"classes": classes,
		},
	})

}

// クラス作成
func RegisterClassHandler(c *gin.Context) {
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

	//構造体に値をバインド
	var bClass model.Class
	if err := c.ShouldBindJSON(&bClass); err != nil {
		fmt.Print("バインド失敗")
		// エラーログ
		return
	}

	// 登録処理を投げてなんかいろいろもらう
	class, err := ClassService.PermissionCheckedClassCreation(idAdjusted, bClass)
	if err != nil {
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
			case common.ErrTypeMaxAttemptsReached: // 最大試行数を超えた
				logging.ErrorLog("Max attempts reached", err)
			}
		} else {
			// エラーログ
			logging.ErrorLog("Class creation was not possible due to other problems.", err)
		}
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": class,
	})
}

func GenerateInviteCodeHandler(c *gin.Context) {
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

	// クラスUUIDを取得
	classUuid := c.Param("class_uuid")

	// 招待コード登録します
	class, err := ClassService.PermissionCheckedRefreshInviteCode(idAdjusted, classUuid)
	if err != nil {
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
			case common.ErrTypeNoResourceExist: // リソースがない
				logging.ErrorLog("The resource does not exist", err)
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			case common.ErrTypeMaxAttemptsReached: // 最大試行数を超えた
				logging.ErrorLog("Max attempts reached", err)
			}
		} else {
			// エラーログ
			logging.ErrorLog("Class creation was not possible due to other problems.", err)
		}
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": class,
	})
}
