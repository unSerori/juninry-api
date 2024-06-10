package route

import (
	"juninry-api/controller"
	"juninry-api/logging"
	"juninry-api/middleware"

	"github.com/gin-gonic/gin"
)

func routing(engine *gin.Engine) *gin.Engine {
	// endpoints
	// root page
	engine.GET("/", controller.ShowRootPage)
	// json test
	engine.GET("/test/json", controller.TestJson) // /test

	// endpoints group
	// ver1グループ
	v1 := engine.Group("/v1")
	{
		// /v1/test
		v1.GET("/test/cfmreq", controller.CfmReq)

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

}

// エンジンを作成して返す
func SetupRouter() (*gin.Engine, error) {
	engine := gin.Default() // エンジンを作成

	routing()

	// テンプレートと静的ファイルを読み込む
	engine.LoadHTMLGlob("view/views/*.html")
	engine.Static("/styles", "./view/views/styles") // クライアントがアクセスするURL, サーバ上のパス
	engine.Static("/scripts", "./view/views/scripts")
	logging.SuccessLog("Routing completed, start the server.")

	return engine, nil // router設定されたengineを返す。
}
