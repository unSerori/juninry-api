// 構造体インスタンスを用いた依存性注入をライブラリで管理
package route

import (
	"juninry-api/application"
	"juninry-api/domain"
	"juninry-api/infrastructure"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/presentation"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"xorm.io/xorm"
)

// エンティティを追加したら、各層のファクトリー関数をprovidersスライスに追加し、ハンドラーの構造体をHandlersに追加

// 依存設定を一括で行うための構造体  // エンティティを追加したら
type Handlers struct {
	dig.In // 継承

	// ハンドラーの構造体 // これをrouter設定側で使えば依存関係をながなが書かなくていい  //　例: v2.POST("/auth/users/t_materials/register", middleware.MidAuthToken(), handlers.TeachingMaterialHandler.RegisterTMHandler) // presentation.NewTeachingMaterialHandler(application.NewTeachingMaterialService(domain.NewTeachingMaterialLogic(), infrastructure.NewTeachingMaterialFileOperatePersistence())).RegisterTMHandler
	TeachingMaterialHandler *presentation.TeachingMaterialHandler
}

// 依存関係を作成
func BuildContainer() *dig.Container {
	// コンテナを作成
	container := dig.New()

	// DBの初期化をコンテナに渡し、依存関係を登録
	container.Provide(
		func() *xorm.Engine {
			db, err := model.InitDB() // router設定されたengineを無名関数でラップしたものを受け取り、ルーティングを登録
			if err != nil {
				logging.ErrorLog("Couldnt receive router engine.", err) // エラー内容を出力し早期リターン
				panic(err)
			}
			return db
		},
	)

	// 登録する依存関係を書く
	providers := []interface{}{
		// 教材登録
		// model.InitDB,
		domain.NewTeachingMaterialLogic,
		infrastructure.NewTeachingMaterialPersistence,
		application.NewTeachingMaterialService,
		presentation.NewTeachingMaterialHandler,
	}

	// スライスから各項目の依存関係を登録し、エラーハンドリング
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			logging.ErrorLog("Dependency registration failed.", nil)
			panic(err)
		}
	}

	// ルーティング設定をコンテナに渡し、依存関係を登録
	container.Provide(
		func(handlers Handlers) *gin.Engine {
			router, err := SetupRouter(handlers) // router設定されたengineを無名関数でラップしたものを受け取り、ルーティングを登録
			if err != nil {
				logging.ErrorLog("Couldnt receive router engine.", err) // エラー内容を出力し早期リターン
				panic(err)
			}
			return router
		},
	)

	return container
}
