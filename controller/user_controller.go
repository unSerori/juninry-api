package controller

import (
	"errors"
	common "juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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
		// 処理で発生したエラーのうちDB関連のエラーのみ
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // 第一引数のerrが第二引数の型にキャスト可能ならキャストしてtrue
			// 本処理時のエラーごとに処理(:DBエラー)
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反
				// エラーログ
				logging.ErrorLog("There is already a user with the same primary key. Uniqueness constraint violation.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default:
				// エラーログ
				logging.ErrorLog("New user registration was not possible due to other DB problems.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		}
		// 処理で発生したエラーのうちDB関連でないもの
		var serviceErr *common.CustomErr
		if errors.As(err, &serviceErr) {
			// 本処理時のエラーごとに処理(:DBエラー以外)
			switch serviceErr.Type {
			case common.ErrTypeHashingPassFailed: // ハッシュ化に失敗
				// エラーログ
				logging.ErrorLog("Failure to hash passwords.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case common.ErrTypeGenTokenFailed: // トークンの作成に失敗
				// エラーログ
				logging.ErrorLog("Failed to generate token.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default:
				// エラーログ
				logging.ErrorLog("New user registration was not possible due to other problems.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			}
		}
		return // エラーレスポンス後に終了
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
	// st := reflect.TypeOf(bUser)  // 型を取得
	// sv := reflect.ValueOf(bUser) // 値を取得
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
			case common.ErrTypeNoResourceExist: // ユーザーが見つからなかった,
				// エラーログ
				logging.ErrorLog("Bad Request.", err)
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // 500番
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
