// ユーザーインターフェース(:リクエストの受け取りとレスポンスの返却)のハンドラー

package presentation

import (
	"errors"
	"juninry-api/application"
	"juninry-api/common/custom"
	"juninry-api/common/responder"
	"juninry-api/dto/requests"
	"juninry-api/dto/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HardHandler struct {
	s *application.HardService
}

// ファクトリー関数
func NewHardHandler(s *application.HardService) *HardHandler {
	return &HardHandler{
		s: s,
	}
}

// ハンドラー

// ハードデバイスの初期設定ハンドラー
func (h *HardHandler) InitHardHandler(ctx *gin.Context) {
	// 構造体にマッピング
	var bReq requests.InitHard // 構造体のインスタンス
	if err := ctx.ShouldBindJSON(&bReq); err != nil {
		responder.SendFailedBindJSON(ctx, err)
		return
	}

	// サービス処理と失敗レスポンス
	id, err := h.s.InitHardService(bReq)
	if err != nil { // エラーハンドル
		// カスタムエラーを仕分ける
		var customErr *custom.CustomErr
		if errors.As(err, &customErr) { // errをcustomErrにアサーションできたらtrue
			switch customErr.Type { // アサーション後のエラータイプで判定 400番台など
			case custom.ErrTypeUniqueConstraintViolation: // 一意性制約違反
				responder.SendFailedService(ctx, http.StatusConflict, err)
			case custom.ErrTypeUnexpectedSetPoints: // 予期せぬ設定値
				responder.SendFailedService(ctx, http.StatusBadRequest, err)
			default: // カスタムエラーの仕分けにぬけがある可能性がある
				responder.SendFailedServiceDefault(ctx, http.StatusInternalServerError, customErr, err)
			}
		} else { // カスタムエラー以外の処理エラー
			responder.SendFailedService(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// 値の整形
	var bRes responses.InitBox = responses.InitBox{ // by dto.responses
		HardwareUuid: id,
	}

	// 成功レスポンス
	responder.SendSuccess(ctx, http.StatusCreated, bRes)
}
