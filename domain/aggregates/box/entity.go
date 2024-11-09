// エンティティの定義

package box

import (
	box "juninry-api/domain/aggregates/box/vo"
	"juninry-api/domain/vo"
)

// エンティティと属性の構造体
type Box struct {
	Id           vo.UUID
	OuchiUuid    vo.UUID
	Status       box.Status
	DepositPoint int
	// これらはrewardが持つもの
	// Name         string
	// RewardPoint  int
	// RewardUuid   vo.UUID
}
