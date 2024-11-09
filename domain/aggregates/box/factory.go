// エンティティのファクトリー関数

package box

import (
	"juninry-api/common/custom"
	box "juninry-api/domain/aggregates/box/vo"
	"juninry-api/domain/vo"
)

// エンティティのファクトリー関数
func NewBox(ouchiId string, status string, depositP int) (*Box, error) {
	// 1. バリデーションに必要な値を作成

	// IDを生成
	boxId, err := vo.NewUUID("")
	if err != nil {
		return nil, err
	}

	// 2. 引数として渡された値のバリデーション

	// おうちID生成
	ouchiUuid, err := vo.NewUUID(ouchiId)
	if err != nil {
		return nil, err
	}
	// // 実在するOuchiUuidか判断
	// if !isRealId(ouchiUuid) {
	// 	logging.ErrorLog("box.NewBox validation.", errors.New("the required character set is not covered"))
	// 	return nil, custom.NewErr(custom.ErrTypeIllegalValue)
	// }

	// ステータスを値オブジェクトとして生成
	var newStatus *box.Status
	if status == "" { // 初期化
		newStatus, err = box.NewStatus("none")
		if err != nil {
			return nil, custom.NewErr(custom.ErrTypeIllegalValue)
		}
	} else {
		newStatus, err = box.NewStatus(status)
		if err != nil {
			return nil, custom.NewErr(custom.ErrTypeIllegalValue)
		}
	}

	// 3. 引数から直接得れるものではないが、エンティティが保持すべき値の生成

	// None

	// 4. 構造体の値として設定しインスタンスを生成、呼び出し元に返却
	return &Box{
		Id:           *boxId,
		OuchiUuid:    *ouchiUuid,
		Status:       *newStatus,
		DepositPoint: depositP,
	}, nil
}

// ビジネスロジック

// // 実在するOuchiUuidか判断
// func isRealId(uuid *vo.UUID) bool {

// 	return true
// }
