package route

import (
	"juninry-api/controller"
	"github.com/gin-gonic/gin"
)

func GetRouter() (*gin.Engine, error) {
	engine := gin.Default() // エンジンを作成

	// endpoints
	// root page
	engine.GET("/", controller.ShowRootPage)
	// json test

	// endpoints group
	// ver1グループ
	v1 := engine.Group("/v1")
	{
		// usersグループ
		users := v1.Group("/users")
		{
			users.POST("/user", controller.RegisterUserHandler)
		}

		// noticesグループ
		notices := v1.Group("/notices")
		{
			//確認できるようにGETにしてる
			notices.GET("/notices",controller.Get_Notices_Handler)
		}
	}

	return engine, nil // router設定されたengineを返す。
}
