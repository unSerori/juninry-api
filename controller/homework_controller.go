package controller

import (
	"errors"
	"fmt"
	"juninry-api/common/logging"
	"juninry-api/model"
	"juninry-api/service"
	"juninry-api/utility/custom"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var homeworkService = service.HomeworkService{}

// 課題の提出履歴を取得
func GetHomeworkRecordHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// リクエストパラメータを取得
	targetMonth := c.Query("targetMonth")
	if targetMonth == "" {
		// エラーログ
		logging.ErrorLog("Don't have targetMonth.", nil)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	//　文字列日付くんを元に戻してあげる
	targetTime, err := time.Parse("2006-01-02 15:04:05.000Z", targetMonth)
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failure to parse date.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	//問い合わせ処理と失敗レスポンス
	homeworkList, err := homeworkService.GetHomeworkRecord(idAdjusted, targetTime)
	if err != nil { //エラーハンドル
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) {
			switch customErr.Type {
			case custom.ErrTypePermissionDenied:
				// エラーログ
				logging.ErrorLog("Permission denied.", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return

			default: // カスタムエラーの仕分けにぬけがある可能性がある
				// エラーログ
				logging.WarningLog("There may be omissions in the CustomErr sorting.", fmt.Sprintf("{customErr.Type: %v, err: %v}", customErr.Type, err))
				// レスポンス
				resStatusCode := http.StatusBadRequest
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
				return
			}
		} else {
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

	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": homeworkList,
	})
}

// 課題全件取得
func FindHomeworkHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	//問い合わせ処理と失敗レスポンス
	homeworkList, err := homeworkService.FindHomework(idAdjusted)
	if err != nil { //エラーハンドル

		var customErr *custom.CustomErr
		if errors.As(err, &customErr) {
			// エラーログ
			logging.ErrorLog(customErr.Error(), nil)
			// レスポンス
			resStatusCode := http.StatusForbidden
			c.JSON(resStatusCode, gin.H{
				"srvResMsg":  http.StatusText(resStatusCode),
				"srvResData": gin.H{},
			})
			return
		}
		// エラーログ
		logging.ErrorLog("SQL query failed.", err)
		//レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get homework list.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": homeworkList,
	})
}

// 次の日の課題を取得
func FindNextdayHomeworkHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	//問い合わせ処理と失敗レスポンス
	homeworkList, err := homeworkService.FindClassHomework(idAdjusted)
	if err != nil { //エラーハンドル
		// エラーログ
		logging.ErrorLog("SQL query failed.", err)
		//レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get homework list.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": homeworkList,
	})
}

// 宿題提出
func SubmitHomeworkHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// form fields 構造体にマッピング
	var bHW *model.HomeworkSubmission    // 構造体のインスタンス
	if err := c.Bind(&bHW); err != nil { // フォームフィールドの直接取得  hwId := c.PostForm("homeworkUUID")
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
	// 構造体にidを追加
	bHW.UserUuid = idAdjusted

	// form files取得
	form, err := c.MultipartForm() // フォームを取得
	if err != nil {
		// エラーログ
		logging.ErrorLog("Failed to retrieve image request.", err)
		// レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 依存性注入
	// fileUploader := &dip.GinContextWrapper{C: c} // サービス層で使えるように、依存性をラッパー構造体のインスタンスとして作成

	// 提出記録処理と失敗レスポンス
	err = homeworkService.SubmitHomework(bHW, form) // 依存性を渡す
	if err != nil {                                 // エラーハンドル
		logging.ErrorLog("Service Error.", err)
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeFileSizeTooLarge: // 画像がでかすぎる
				// エラーログ
				logging.ErrorLog("Payload Too Large.", err)
				// レスポンス
				resStatusCode := http.StatusRequestEntityTooLarge
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypeInvalidFileFormat: // 画像形式が不正
				// エラーログ
				logging.ErrorLog("Unsupported Media Type.", err)
				// レスポンス
				resStatusCode := http.StatusUnsupportedMediaType
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
	logging.SuccessLog("Successful submission homework.")
	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": gin.H{},
	})
}

// 宿題登録
func RegisterHWHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// 構造体にマッピング
	var bHW service.BindRegisterHW // 構造体のインスタンス
	if err := c.ShouldBindJSON(&bHW); err != nil {
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
	hwId, err := homeworkService.RegisterHWService(bHW, idAdjusted)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeUniqueConstraintViolation: // 一意性制約違反
				// エラーログ
				logging.ErrorLog("Conflict.", err)
				// レスポンス
				resStatusCode := http.StatusConflict
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
			"homeworkUUID": hwId,
		},
	})
}

// 特定の宿題に対する任意のユーザーの提出状況と宿題の詳細情報を取得するエンドポイント
func GetHWInfoHandler(c *gin.Context) {
	// ユーザーを特定する(ctxに保存されているidを取ってくる)
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// パスパラhomework_uuidの取得
	homeworkUuid := c.Param("homework_uuid")
	// クエパラuser_uuidの取得
	userUuid := c.Query("user_uuid")

	fmt.Printf("idAdjusted: %v\n", idAdjusted)
	fmt.Printf("homeworkUuid: %v\n", homeworkUuid)
	if userUuid == "" {
		fmt.Println("emp string")
	} else {
		fmt.Printf("userUuid: %v\n", userUuid)
	}

	// 取得処理と失敗レスポンス
	resData, err := homeworkService.GetHWInfoService(homeworkUuid, idAdjusted, userUuid)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // リソースがなく見つからない
				// エラーログ
				logging.ErrorLog("Not Found.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypePermissionDenied: // 所属していないクラスのお知らせを取得しようとしているね
				// エラーログ
				logging.ErrorLog("Forbidden.", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
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

	// 成功ログ
	logging.SuccessLog("Successful noticeDetail get.")
	// レスポンス(StatusOK　成功200番)
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg":  http.StatusText(resStatusCode),
		"srvResData": resData,
	})
}

// 教材データを取得
func GetMaterialDataHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	classUuid := c.Param("classUuid")

	//問い合わせ処理と失敗レスポンス
	materialList, err := homeworkService.GetTeachingMaterialData(idAdjusted, classUuid)
	if err != nil { //エラーハンドル
		// エラーログ
		logging.ErrorLog("SQL query failed.", err)
		//レスポンス
		resStatusCode := http.StatusBadRequest
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get material list.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"teachingItems": materialList,
		},
	})
}

// 特定の提出済み宿題の画像を取得する
func FetchSubmittedHwImageHandler(c *gin.Context) {
	// ユーザーを特定する(ctxに保存されているidを取ってくる)
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// パスパラの取得
	homeworkUuid := c.Param("homework_uuid")    // クラス
	imageFileName := c.Param("image_file_name") // 画像パス

	// パス作成処理と失敗レスポンス
	filePath, err := homeworkService.FetchSubmittedHwImageService(idAdjusted, homeworkUuid, imageFileName)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeNoResourceExist: // リソースがなく見つからない
				// エラーログ
				logging.ErrorLog("Not Found.", err)
				// レスポンス
				resStatusCode := http.StatusNotFound
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case custom.ErrTypePermissionDenied: // 画像へのアクセスなし
				// エラーログ
				logging.ErrorLog("Forbidden.", err)
				// レスポンス
				resStatusCode := http.StatusForbidden
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
	logging.SuccessLog("Successful image acquisition.")
	// 画像レスポンス
	c.File(filePath)
}

// 教師が特定の宿題に対するその宿題が配られたクラスの生徒の進捗一覧を取得
func GetStudentsHomeworkProgressHandler(c *gin.Context) {
	// ユーザーを特定する
	id, _ := c.Get("id")
	idAdjusted := id.(string) // アサーション

	// パスパラの取得
	hwId := c.Param("homework_uuid") // hwId

	// 一覧取得処理と失敗レスポンス
	studentSubmissionInfoSlice, err := homeworkService.GetStudentsHomeworkProgressService(idAdjusted, hwId)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeUniqueConstraintViolation: // 一意性制約違反
				// エラーログ
				logging.ErrorLog("Conflict.", err)
				// レスポンス
				resStatusCode := http.StatusConflict
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
			"progress": studentSubmissionInfoSlice,
		},
	})
}
