// 構造体インスタンスを用いた依存性注入をライブラリで管理
package di

import (
	"juninry-api/logging"
	"juninry-api/route"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// 依存関係を作成
func BuildContainer() *dig.Container {
	// コンテナを作成
	container := dig.New()

	// 登録する依存関係を書く
	providers := []interface{}{
		// テスト

		// router設定されたengineを無名関数でラップしたものを受け取り、ルーティングを登録
		func() *gin.Engine {
			router, err := route.SetupRouter()
			if err != nil {
				logging.ErrorLog("Couldnt receive router engine.", err) // エラー内容を出力し早期リターン
				panic(err)
			}
			return router
		},
	}

	// スライスから各項目の依存関係を登録し、エラーハンドリング
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			logging.ErrorLog("Couldnt receive router engine.", nil)
			panic(err)
		}
	}

	return container
}
