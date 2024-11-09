// ORM関連のリポジトリインターフェース

package box

import "juninry-api/model"

// リポジトリインターフェースの定義
type OrmRepoI interface {
	// ouchiUuidが存在するか確認
	CheckOuchiUuid(ouchiUuid string) error
	// ボックスの新規登録
	AddBox(record model.Box) error
}
