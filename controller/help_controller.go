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

var helpService = service.HelpService{} // サービスの実体を作る。

// おてつだいを取得
func GetHelpsHandler(c *gin.Context) {
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
	// おてつだい
	helps, err := helpService.GetHelps(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
	}
	// 成功ログ
	logging.SuccessLog("Successful user get.")
	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		// "srvResCode": 1001,
		// "srvResMsg":  "Successful user get.",
		"srvResData": gin.H{
			"rewardData": helps,
		},
	})
}

// おてつだい追加
func CreateHelpHandler(c *gin.Context) {
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

	// 構造体に値をバインド
	var bHelp model.Help
	if err := c.ShouldBindBodyWithJSON(&bHelp); err != nil {
		fmt.Print("バインド失敗")
		return
	}

	// おてつだいを作成
	helps, err := helpService.CreateHelp(idAdjusted, bHelp)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
	}
	// 成功ログ
	logging.SuccessLog("Successful user get.")
	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		// "srvResCode": 1001,
		// "srvResMsg":  "Successful user get.",
		"srvResData": gin.H{
			"rewardData": helps,
		},
	})
}

// おてつだいを消化
func HelpSubmittionHandler(c *gin.Context) {
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

	// 構造体に値をバインド
	var bHelps model.HelpSubmittion
	if err := c.ShouldBindBodyWithJSON(&bHelps); err != nil {
		fmt.Print("バインド失敗")
		return
	}

	// おてつだいを交換
	helps, err := helpService.HelpDigestion(idAdjusted, bHelps)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
	}
	// 成功ログ
	logging.SuccessLog("Successful user get.")
	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		// "srvResCode": 1001,
		// "srvResMsg":  "Successful user get.",
		"srvResData": gin.H{
			"rewardData": helps,
		},
	})
}