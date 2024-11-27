// ロジックから呼び出される具体的な永続化処理
package infrastructure

import (
	"juninry-api/common/custom"
	"juninry-api/common/logging"
	ri "juninry-api/domain/aggregates/box/ri"
	"juninry-api/model"

	"xorm.io/xorm"
)

// 技術実装の構造体
type BoxOrmRepoImpl struct {
	db *xorm.Engine
}

// ファクトリー関数
func NewBoxOrmRepoImpl(db *xorm.Engine) ri.OrmRepoI {
	return &BoxOrmRepoImpl{db: db}
}

// 実装

// ouchiUuidが存在するか確認
func (i *BoxOrmRepoImpl) CheckOuchiUuid(ouchiUuid string) error {
	var box model.Box // 取得したデータをマッピングする構造体

	// 取得
	isFound, err := i.db.Where("ouchi_uuid = ?", ouchiUuid).Get(&box)
	// エラーハンドル
	if err != nil {
		logging.ErrorLog("Error when CheckOuchiUuid err.", err)
		return err
	}
	if !isFound {
		logging.ErrorLog("Error when CheckOuchiUuid isFound.", err)
		return custom.NewErr(custom.ErrTypeNoFoundR)
	}

	return nil
}

// ボックスの新規登録
func (i *BoxOrmRepoImpl) AddBox(record model.Box) error {

	return nil
}
