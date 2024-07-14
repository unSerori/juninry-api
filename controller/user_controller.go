package controller

import (
	"errors"
	"fmt"
	common "juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userService = service.UserService{} // サービスの実体を作る。

// ユーザ作成
func RegisterUserHandler(c *gin.Context) {
	// 構造体にマッピング
	var bUser model.User // 構造体のインスタンス
	if err := c.ShouldBindJSON(&bUser); err != nil {
		// エラーログ
		logging.ErrorLog("Failure to bind request.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 登録処理と失敗レスポンス
	token, err := userService.RegisterUser(bUser)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeUniqueConstraintViolation: //
				// エラーログ
				logging.ErrorLog("Bad Request.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
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
	logging.SuccessLog("Successful user registration.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"authenticationToken": token,
		},
	})
}

// ユーザ情報取得
func GetUserHandler(c *gin.Context) {
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

	// 構造体にマッピング
	var bUser model.User // 構造体のインスタンス

	// ユーザー情報の取得
	// bUser, err := userService.GetUser(userData.UserUuid);
	bUser, err := userService.GetUser(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failure to get user.", err)
		// レスポンス
		c.JSON(http.StatusBadRequest, gin.H{
			// "srvResCode": 7001,
			// "srvResMsg":  "Failure to getuser.",
			"srvResData": gin.H{},
		})
		return
	}

	// 成功ログ
	logging.SuccessLog("Successful user get.")
	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		// "srvResCode": 1001,
		// "srvResMsg":  "Successful user get.",
		"srvResData": gin.H{
			"userData": bUser,
		},
	})
}

// login
func LoginHandler(c *gin.Context) {
	// 構造体にマッピング
	var bUser model.User // 構造体のインスタンス
	if err := c.ShouldBindJSON(&bUser); err != nil {
		// エラーログ
		logging.ErrorLog("Failure to bind request.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	// // 構造体の中身をチェック
	// st := reflect.TypeOf(bUser)  // 型を取得 構造体がポインタなら*ptrを渡す
	// sv := reflect.ValueOf(bUser) // 値を取得 構造体がポインタなら*ptrを渡す
	// // 構造体のフィールド数だけループ
	// for i := 0; i < st.NumField(); i++ {
	// 	fieldName := st.Field(i).Name                             // フィールド名を取得
	// 	fieldValue := sv.Field(i)                                 // フィールドの値を取得
	// 	fmt.Printf("%s: %v\n", fieldName, fieldValue.Interface()) // フィールド名と値を出力
	// }

	// ログイン処理と失敗レスポンス
	token, err := userService.LoginUser(bUser)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeNoResourceExist, common.ErrTypePassMismatch: // ユーザーが見つからなかった, パスワードが不一致
				// エラーログ
				logging.ErrorLog("Bad Request.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
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

	// 成功レスポンス 200番
	// 成功ログ
	logging.SuccessLog("Successful user login.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"authenticationToken": token,
		},
	})
}
