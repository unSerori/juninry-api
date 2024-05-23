package route

import (
	"juninry-api/controller"

	"github.com/gin-gonic/gin"
)

func GetRouter() (*gin.Engine, error) {
	engine := gin.Default() // エンジンを作成

	// endpoint
	// root page
	engine.GET("/", controller.ShowRootPage)
	// json test

	// endpoints group
	// apiのグループ

	return engine, nil // router設定されたengineを返す。
}
