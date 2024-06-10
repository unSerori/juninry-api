package route

import (
	"juninry-api/controller"
	"juninry-api/middleware"

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
		// /v1/test
		v1.GET("test", controller.TestJson)

		// usersグループ
		users := v1.Group("/users")
		{
			// ユーザー新規登録 /v1/users/register
			users.POST("/register", controller.RegisterUserHandler)
		}

		// authグループ
		auth := v1.Group("/auth", middleware.MidLog())
		auth.Use(middleware.MidAuthToken()) // 認証ミドルウェア適用
		{
			// /v1/auth/test/cfmreq
			auth.GET("/test/cfmreq", controller.CfmReq)

			// usersグループ
			users := auth.Group("/users")
			{
				// homeworkグループ
				homework := users.Group("/homework")
				{
					// /v1/auth/users/homework/upcoming
					homework.GET("/test", controller.CfmReq)
				}

				// noticeグループ
				notice := users.Group("/notice")
				{
					// /v1/auth/users/notice/notices
					notice.GET("/notices", controller.CfmReq)

					// /v1/auth/users/notice/{notice_uuid}
					notice.GET("/:notice_uuid", controller.TestJson) // おしらせ詳細をとる // コントローラで取り出すときは noticeUuid := c.Param("notice_uuid")
				}

			}
		}
	}

	return engine, nil // router設定されたengineを返す。
}
