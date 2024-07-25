// リクエストの取得し、ビジネスロック(アプリケーション層の関数)を呼び出し、レスポンスを返す

package presentation

import (
	"errors"
	"fmt"
	"juninry-api/application"
	"juninry-api/common"
	"juninry-api/domain"
	"juninry-api/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 依存関係のための構造体
type TeachingMaterialHandler struct {
	s *application.TeachingMaterialService // 依存先層の構造体 依存先層のポインタ型
}

// 依存関係のためのファクトリー関数
func NewTeachingMaterialHandler(s *application.TeachingMaterialService) *TeachingMaterialHandler { // 依存先層の構造体 依存先層のポインタ型
	return &TeachingMaterialHandler{s: s} // &Handler{依存先層の構造体: 依存先層の構造体}
}

// 教科作成用エンドポイント
func (h TeachingMaterialHandler) RegisterTMHandler(c *gin.Context) {
	// ユーザーを特定する
	id, exists := c.Get("id")
	if !exists { // idがcに保存されていない。
		// エラーログ
		logging.ErrorLog("Internal Server Error.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}
	idAdjusted := id.(string) // アサーション

	// form fields 構造体にマッピング
	var bTM domain.TeachingMaterial // 構造体のインスタンス
	// bTM := new(domain.TeachingMaterial) && c.Bind(bTM)
	if err := c.Bind(&bTM); err != nil { // フォームフィールドの直接取得 => hwId := c.PostForm("fields")
		// エラーログ
		logging.ErrorLog("Internal Server Error.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// form files取得 // file, _ := c.FormFile("file")でも最初の1ファイル目が取得できるが、今後の処理・エラー拡張性のため全部取る
	form, err := c.MultipartForm() // フォームを取得
	if err != nil {
		// エラーログ
		logging.ErrorLog("Internal Server Error.", nil)
		// レスポンス
		resStatusCode := http.StatusInternalServerError
		c.JSON(resStatusCode, gin.H{
			"srvResMsg":  http.StatusText(resStatusCode),
			"srvResData": gin.H{},
		})
		return
	}

	// 新規教材作成処理と失敗レスポンス
	tmId, err := h.s.RegisterTMService(idAdjusted, bTM, form) // idAdjusted, bTM, form
	if err != nil {                                           // エラーハンドル
		logging.ErrorLog("Service Error.", err)
		// カスタムエラーを仕分ける
		var customErr *common.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case common.ErrTypeFileSizeTooLarge: // 画像がでかすぎる
				// エラーログ
				logging.ErrorLog("Payload Too Large.", err)
				// レスポンス
				resStatusCode := http.StatusRequestEntityTooLarge
				c.JSON(resStatusCode, gin.H{
					"srvResMsg":  http.StatusText(resStatusCode),
					"srvResData": gin.H{},
				})
			case common.ErrTypeInvalidFileFormat: // 画像形式が不正
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
	logging.SuccessLog("Successful submission teaching material.")
	// レスポンス
	resStatusCode := http.StatusCreated
	c.JSON(resStatusCode, gin.H{
		"srvResMsg": http.StatusText(resStatusCode),
		"srvResData": gin.H{
			"teachingMaterialUuid": tmId,
		},
	})
}
