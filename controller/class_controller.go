package controller

import (
	"errors"
	"fmt"
	"juninry-api/common/custom"
	"juninry-api/common/logging"
	"juninry-api/model"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ClassService = service.ClassService{}

// クラス一覧取得
func GetAllClassesHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// サービスに投げるよ
	classes, err := ClassService.GetClassList(idAdjusted)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to get class list.", err)
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // カスタムエラーの場合
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // // お家に子供いないよエラー
				// エラーログ
				logging.ErrorLog("Bad Request.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{ // お家に子供いないよエラー
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
	id, _ := c.Get("id")
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
		var serviceErr *custom.CustomErr
		if errors.As(err, &serviceErr) { // カスタムエラーの場合
			switch serviceErr.Type {
			case custom.ErrTypePermissionDenied: // 権限を持っていない
				logging.ErrorLog("Do not have the necessary permissions", err)
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypeMaxAttemptsReached: // 最大試行数を超えた
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
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", serviceErr.Type, err))
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

	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": class,
	})
}

// ユーザーIDから参加しているクラスを取得し、生徒一覧を返す
func GetClasssmaitesHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// idからクラスメイトの情報を取得
	classmates, err := ClassService.GetClassMates(idAdjusted)
	// エラーハンドル
	if err != nil {
		//カスタムエラーを分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // ユーザーが見つからなかった, パスワードが不一致
				// エラーログ
				logging.ErrorLog("Not Found.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
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
		} else {
			// エラーログ
			logging.ErrorLog("Failure to get user.", err)
			// レスポンス
			resStatusCode := http.StatusBadRequest
			c.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})

		}

		return
	}

	// 成功ログ
	logging.SuccessLog("Successful users get.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": classmates,
	})
}

func GenerateInviteCodeHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// クラスUUIDを取得
	classUuid := c.Param("class_uuid")

	fmt.Printf("classUuid: %v\n", classUuid)

	// 招待コード登録します
	class, err := ClassService.PermissionCheckedRefreshInviteCode(idAdjusted, classUuid)
	if err != nil {
		var serviceErr *custom.CustomErr
		if errors.As(err, &serviceErr) { // カスタムエラーの場合
			switch serviceErr.Type {
			case custom.ErrTypePermissionDenied: // 権限を持っていない
				logging.ErrorLog("Do not have the necessary permissions", err)
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypeNoResourceExist: // リソースがない
				logging.ErrorLog("The resource does not exist", err)
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", serviceErr.Type, err))
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

	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": class,
	})
}

func JoinClassHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// クラスUUIDを取得
	inviteCode := c.Param("invite_code")

	// 出席番号受け取りマン
	var studentNumberJSON struct {
		StudentNumber *int `json:"studentNumber"`
	}

	// JSONをバインド
	err := c.ShouldBind(&studentNumberJSON)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Bind error.", err)
		// TODO:JSONがちゃんと送られてきていない場合と、ミスってる場合の切り分けができていないけど手段なくね
	}

	// クラスに参加
	className, err := ClassService.PermissionCheckedJoinClass(idAdjusted, inviteCode, studentNumberJSON.StudentNumber)
	if err != nil {
		var serviceErr *custom.CustomErr
		if errors.As(err, &serviceErr) { // カスタムエラーの場合
			switch serviceErr.Type {
			case custom.ErrTypeUniqueConstraintViolation: // 一意性制約違反
				// エラーログ
				logging.ErrorLog("Conflict.", err)
				// レスポンス
				resStatusCode := http.StatusConflict
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypePermissionDenied: // 権限を持っていない
				logging.ErrorLog("Do not have the necessary permissions", err)
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypeNoResourceExist: // 招待コード違います
				logging.ErrorLog("The resource does not exist", err)
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", serviceErr.Type, err))
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

	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"className": className,
		},
	})

}
