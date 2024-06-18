package route

import (
	"juninry-api/controller"
	"juninry-api/logging"
	"juninry-api/middleware"

	"github.com/gin-gonic/gin"
)

func routing(engine *gin.Engine) {
	// MidLog all
	engine.Use(middleware.LoggingMid())

	// endpoints
	// root page
	engine.GET("/", controller.ShowRootPage) // /
	// json test
	engine.GET("/test/json", controller.TestJson) // /test

	// endpoints group
	// ver1グループ
	v1 := engine.Group("/v1")
	{
		// リクエストを鯖側で確かめるテスト用エンドポイント
		v1.GET("/test/cfmreq", controller.CfmReq) // /v1/test

		// usersグループ
		users := v1.Group("/users")
		{
			// ユーザー新規登録
			users.POST("/register", controller.RegisterUserHandler) // /v1/users/register

			// ユーザーログイン
			users.POST("/login", controller.LoginHandler) // /v1/users/login
		}

		// authグループ 認証ミドルウェア適用
		auth := v1.Group("/auth", middleware.MidAuthToken())
		{
			// 認証グループで、認証ができるかを確認するテスト用エンドポイント
			auth.GET("/test/cfmreq", controller.CfmReq) // /v1/auth/test/cfmreq

			// usersグループ
			users := auth.Group("/users")
			{
				// ユーザー自身のプロフィールを取得
				users.GET("/user", controller.GetUserHandler) // /v1/auth/auth/users/user

				// homeworksグループ
				homeworks := users.Group("/homework")
				{
					// 期限がある課題一覧を取得
					homeworks.GET("/upcoming", controller.FindHomeworkHandler) // /v1/auth/users/homework/upcoming
				}

				// noticeグループ
				notices := users.Group("/notice")
				{
					// 自分の所属するクラスのおしらせ一覧をとる
					notices.GET("/notices", controller.CfmReq) // /v1/auth/users/notice/notices

					// おしらせ詳細をとる // コントローラで取り出すときは noticeUuid := c.Param("notice_uuid")
					notices.GET("/:notice_uuid", controller.TestJson) // /v1/auth/users/notice/{notice_uuid}
				}
			}
		}
	}

}

// ファイルを設定
func loadingStaticFile(engine *gin.Engine) {
	// テンプレートと静的ファイルを読み込む
	engine.LoadHTMLGlob("view/views/*.html")
	engine.Static("/styles", "./view/views/styles") // クライアントがアクセスするURL, サーバ上のパス
	engine.Static("/scripts", "./view/views/scripts")
	logging.SuccessLog("Routing completed, start the server.")

}

// エンジンを作成して返す
func SetupRouter() (*gin.Engine, error) {
	// エンジンを作成
	engine := gin.Default()

	// ルーティング
	routing(engine)

	// 静的ファイル設定
	loadingStaticFile(engine)

	// router設定されたengineを返す。
	return engine, nil
}
