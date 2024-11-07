package controller

import (
	"errors"
	"fmt"
	"juninry-api/common/logging"
	"juninry-api/model"
	"juninry-api/service"
	"juninry-api/utility/custom"
	"net/http"

	"github.com/gin-gonic/gin"
)

var rewardService = service.RewardService{} // サービスの実体を作る。

// ごほうびを取得
func GetRewardsHandler(c *gin.Context) {
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
	// ご褒美を取得
	rewards, err := rewardService.GetRewards(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
			"rewardData": rewards,
		},
	})
}

func GetBoxRewardsHandler(c *gin.Context) {
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
	// ご褒美を取得

	rewards, err := rewardService.GetBoxRewards(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
	resStatusCode := http.StatusOK
	c.JSON(http.StatusCreated, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"boxes": rewards,
		},
	})
}

// ごほうびを作成
func CreateRewardHandler(c *gin.Context) {
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
	var bReward model.Reward
	if err := c.ShouldBindBodyWithJSON(&bReward); err != nil {
		fmt.Print("バインド失敗")
		return
	}

	// ご褒美を作成
	rewards, err := rewardService.CreateRewards(idAdjusted, bReward)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
			"rewardData": rewards,
		},
	})
}

// ごほうびを削除
func DeleteRewardsHandler(c *gin.Context) {
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

	// reward_UUIDを取得
	rewardUuid := c.Param("reward_uuid")
	fmt.Printf("rewardUUID: %v\n", rewardUuid)

	// ご褒美を削除
	_, err := rewardService.DeleteReward(idAdjusted, rewardUuid)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
			"srvResMsg": "ok!",
		},
	})
}

// ごほうびを交換
func RewardsExchangeHandler(c *gin.Context) {
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
	var bReward model.RewardExchanging
	if err := c.ShouldBindBodyWithJSON(&bReward); err != nil {
		fmt.Print("バインド失敗")
		return
	}

	// ご褒美を交換
	ouchiPoint, err := rewardService.ExchangeReward(idAdjusted, bReward)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
			"ouchiPoint": ouchiPoint,
		},
	})
}

// ごほうびを消化
func RewardDigestionHandler(c *gin.Context) {
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

	// reward_UUIDを取得
	rewardExchangeId := c.Param("rewardExchangeId")
	fmt.Printf("rewardUUID: %v\n", rewardExchangeId)

	// 交換されたご褒美を消化
	_, err := rewardService.RewardDigestion(idAdjusted, rewardExchangeId)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
			"rewardData": "ok!",
		},
	})

}

// 交換されたごほうび一覧を取得
func GetExchangedRewardsHandler(c *gin.Context) {
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

	// 交換されたご褒美を取得
	result, err := rewardService.GetRewardExchanging(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー

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
		"srvResData": gin.H{
			"exchangeData": result,
		},
	})
}

func DepositPointHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")

	idAdjusted := id.(string) // アサーション

	// 追加するポイントを取得
	// リクエストボディのJSONをマップにバインド
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	// jsonにaddPointがあるかを確認
	addPoint, ok := jsonData["addPoint"]
	if !ok || addPoint == nil {
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	addPointInt := int(jsonData["addPoint"].(float64))

	// パスパラメータのハードウェアIDも取得
	hardUuid := c.Param("hardware_uuid")
	if hardUuid == "" {
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 交換されたごほうび一覧を取得
	result, err := rewardService.BoxAddPoint(idAdjusted, addPointInt, hardUuid)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家ないですよエラー
				fmt.Println("No resource exist.")
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return

			case custom.ErrTypeUnforeseenCircumstances: // ポイントがおかしい
				fmt.Println("Unforeseen circumstances.")
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			case custom.ErrTypePermissionDenied: // 権限がありません
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			default: // エラーハンドリング漏れ
				resStatusCode := http.StatusInternalServerError
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
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
			return
		}
	}

	// 成功ログ
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"depositPoint": result,
		},
	})

}
