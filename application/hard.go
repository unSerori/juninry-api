// サービスのユースケースを書く(:処理の流れ)

package application

import (
	"juninry-api/common/custom"
	"juninry-api/domain/aggregates/box"
	ri "juninry-api/domain/aggregates/box/ri"
	"juninry-api/dto/requests"
	"juninry-api/model"
)

// サービスの構造体
type HardService struct {
	r ri.OrmRepoI
}

// ファクトリー関数
func NewHardService(r ri.OrmRepoI) *HardService {
	return &HardService{r: r}
}

// ビジネスロジック

// 宝箱の初期設定サービス
func (s *HardService) InitHardService(req requests.InitHard) (string, error) {
	// 種類によって処理を分岐
	switch req.HardwareTypeId {
	case 1: // 宝箱
		// 必要な値を受け取り初期化
		newBox, err := box.NewBox(
			req.OuchiUuid,
			"",
			0,
		)
		if err != nil {
			return "", err
		}

		// DBに登録
		err = s.r.AddBox(model.Box{
			HardwareUuid: newBox.Id.Value(),
			DepositPoint: newBox.DepositPoint,
			BoxStatus:    newBox.Status.AsInt(),
			OuchiUuid:    newBox.OuchiUuid.Value(),
		})
		if err != nil {
			return "", err
		}

		return newBox.Id.Value(), nil
	default:
		return "", custom.NewErr(custom.ErrTypeUnexpectedSetPoints)
	}
}
