// ハートウェアタイプの値オブジェクト

package box

import "juninry-api/common/custom"

// 値オブジェクトの定義
type Status struct {
	value string
}

// さまざまな値
var statuses = []Status{
	{"none"},
	{"live"},
	{"max"},
	{"maint"},
}

// ファクトリー関数
func NewStatus(value string) (*Status, error) {
	// 1. 引数として渡された値のバリデーション

	// あらかじめ決められた値かどうか
	var isValidStatus bool = false
	for _, status := range statuses {
		if status.value == value {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus { // 見つからなかった
		return nil, custom.NewErr(custom.ErrTypeIllegalValue)
	}

	// 2. 引数値以外のエンティティが保持すべき値の生成

	// None

	// 3. 構造体の値として設定しインスタンスを生成、呼び出し元に返却
	return &Status{
		value: value,
	}, nil

}

// アクセサ
func (vo *Status) Value() string {
	return vo.value
}

// ビジネスロジック
func (vo *Status) AsInt() int {
	// ステータスがあらかじめ定義したスライスのいくつめか探す
	for i, status := range statuses {
		if vo.value == status.value { // 合致したときのindexを返す
			return i
		}
	}

	return -1
}
