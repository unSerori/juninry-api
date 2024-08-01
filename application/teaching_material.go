// サービス(:処理の流れ)

package application

import (
	"juninry-api/domain"
	"mime/multipart"

	"github.com/google/uuid"
)

// アプリケーション層の構造体
type TeachingMaterialService struct {
	// ビジネスロジック
	l *domain.TeachingMaterialLogic
}

// ファクトリー関数でdi
func NewTeachingMaterialService(l *domain.TeachingMaterialLogic) *TeachingMaterialService {
	return &TeachingMaterialService{
		l: l,
	}
}

// 教材追加
func (s *TeachingMaterialService) RegisterTMService(id string, bTM domain.TeachingMaterial, form *multipart.Form) (string, error) {
	// idと所属クラスと教科についての確認
	tm, err := s.l.IntegrityStruct(id, bTM)
	if err != nil {
		return "", err
	}

	// 画像保存
	filepath, err := s.l.ValidateImage(form)
	if err != nil {
		return "", err
	}

	// 教科IDを作成
	submitId, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {
		return "", err
	}
	tm.TeachingMaterialUuid = submitId.String() // 設定

	// 取得したファイル名を構造体に設定して、テーブルにインサート
	tm.TeachingMaterialImageUuid = filepath
	err = s.l.CreateTM(tm)
	if err != nil {
		return "", err
	}

	return tm.TeachingMaterialUuid, nil
}
