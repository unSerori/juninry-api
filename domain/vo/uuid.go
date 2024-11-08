// 共通で使う値オブジェクト

package vo

import (
	"errors"
	"juninry-api/common/custom"
	"juninry-api/common/logging"

	"github.com/google/uuid"
)

// 値オブジェクトの定義
type UUID struct {
	value string
}

// ファクトリー関数
func NewUUID(id string) (*UUID, error) { // UUID生成ルール
	// 引数がない場合は新規UUIDを生成
	// 引数がある場合はバリデーションして値オブジェクトを返す
	if id == "" {
		// 1. 引数として渡された値のバリデーション

		// None

		// 2. 引数値以外のエンティティが保持すべき値の生成

		// 新しいIDの生成
		newId, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		// 3. 構造体の値として設定しインスタンスを生成、呼び出し元に返却
		return &UUID{
			value: newId.String(),
		}, nil
	} else {
		// 1. 引数として渡された値のバリデーション

		// UUIDとして有効かどうか
		if !isValidUUID(id) {
			logging.ErrorLog("user.NewPassword() validation", errors.New("contains double-byte characters"))
			return nil, custom.NewErr(custom.ErrTypeIllegalValue)
		}

		// 2. 引数値以外のエンティティが保持すべき値の生成

		// None

		// 3. 構造体の値として設定しインスタンスを生成、呼び出し元に返却
		return &UUID{
			value: id,
		}, nil
	}
}

// アクセサ
func (vo *UUID) Value() string {
	return vo.value
}

// ビジネスロジック

// UUIDとして有効かどうか
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
