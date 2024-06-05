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
			// /v1/users/register
			users.POST("/register", controller.RegisterUserHandler)
		}

		// authグループ
		auth := v1.Group("/auth")
		{
			// classグループ
			class := auth.Group("/class")
			{
				// /v1/auth/class/notice
				class.GET("/test", controller.TestJson)

				// homeworkグループ
				homework := class.Group("/homework")
				{
					// /v1/auth/class/homework/upcoming
					homework.GET("/test", controller.TestJson)
				}

				// noticeグループ
				notice := class.Group("/notice")
				{
					// /v1/auth/class/notice/{notice_uuid}
					notice.GET("/:notice_uuid", controller.TestJson)
					// コントローラで取り出すときは noticeUuid := c.Param("notice_uuid")
				}
			}
		}
	}

	return engine, nil // router設定されたengineを返す。
}
