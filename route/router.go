package route

import (
	"juninry-api/common"
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
		v1.GET("/test/cfmreq", controller.CfmReq) // /v1/test/cfmreq

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
				homeworks := users.Group("/homeworks")
				{
					// 期限がある課題一覧を取得
					homeworks.GET("/upcoming", controller.FindHomeworkHandler) // /v1/auth/users/homeworks/upcoming

					// 次の日が期限の課題一覧を取得
					homeworks.GET("/nextday", controller.FindNextdayHomeworkHandler) // /v1/auth/users/homeworks/upcoming

					// 宿題の提出
					homeworks.POST("/submit", middleware.LimitReqBodySize(common.LoadReqBodyMaxSize(10485760)), controller.SubmitHomeworkHandler) // /v1/auth/users/homeworks/submit // リクエスト制限のデフォ値は10MB
				}

				// noticeグループ
				notices := users.Group("/notices")
				{
					// 自分の所属するクラスのおしらせ一覧をとる
					notices.GET("/notices", controller.GetAllNoticesHandler) // /v1/auth/users/notices/notices

					// おしらせ詳細をとる // コントローラで取り出すときは noticeUuid := c.Param("notice_uuid")
					notices.GET("/:notice_uuid", controller.GetNoticeDetailHandler) // /v1/auth/users/notice/{notice_uuid}

					//　お知らせ新規登録
					notices.POST("/register", controller.RegisterNoticeHandler) // /v1/auth/users/notices/register

					// お知らせ既読済み処理
					notices.POST("/read/:notice_uuid", controller.NoticeReadHandler) // /v1/auth/users/notices/read/{notice_uuid}

					// 特定のお知らせを既読しているユーザ一覧を取る(エンドポイント名不安。)
					notices.GET("/status/:notice_uuid", controller.GetNoticestatusHandler) // /v1/auth/users/notices/status/{notice_uuid}
				}

				// classesグループ
				classes := users.Group("/classes")
				{

					// クラスに所属する人間たちを返す
					classes.GET("/users", controller.GetClasssmaitesHandler) // /v1/auth/users/classes/users

					// 自分の所属するクラス一覧をとる
					classes.GET("/affiliations", controller.GetAllClassesHandler) // /v1/auth/users/classes/classes

					// クラスを作成する
					classes.POST("/register", middleware.SingleExecutionMiddleware(), controller.RegisterClassHandler) // /v1/auth/users/classes/register

					// 招待コードを更新する
					classes.PUT("/refresh/:class_uuid", controller.GenerateInviteCodeHandler) // /v1/auth/users/classes/invite-code

					// クラスに参加する
					classes.POST("/join/:invite_code", controller.JoinClassHandler)
				}

				// ouchiesグループ
				ouchies := users.Group("/ouchies")
				{
					// おうち作成
					ouchies.POST("/register", middleware.SingleExecutionMiddleware(), controller.RegisterOuchiHandler) // /v1/auth/users/ouchies/register

					// 招待コードの更新
					ouchies.PUT("/refresh/:ouchi_uuid", controller.GenerateOuchiInviteCodeHandler) // /v1/auth/users/ouchies/refresh/{ouchi_uuid}

					// おうちに所属
					ouchies.POST("/join/:invite_code", controller.JoinOuchiHandler) // /v1/auth/users/ouchies/join/{invite_code}

					helps := ouchies.Group("/helps")
					{

						//おてつだいを取得
						helps.GET("/helps", controller.GetHelpsHandler) // /v1/auth/users/ouchies/register

						// おてつだいを追加
						helps.POST("/register", middleware.SingleExecutionMiddleware(), controller.CreateHelpHandler) // /v1/auth/users/ouchies/join/{invite_code}

						// おてつだいを消化
						helps.POST("/submittion", controller.HelpSubmittionHandler) // /v1/auth/users/ouchies/refresh/{ouchi_uuid}
					}

					rewards := ouchies.Group("/rewards")
					{

						// ごほうびを取得
						rewards.GET("/rewards", controller.GetRewardsHandler) // /v1/auth/users/ouchies/register

						// ごほうびを追加
						rewards.POST("/register", middleware.SingleExecutionMiddleware(), controller.CreateRewardHandler) // /v1/auth/users/ouchies/join/{invite_code}

						// ごほうびを交換
						rewards.POST("/exchange", controller.RewardsExchangeHandler) // /v1/auth/users/ouchies/refresh/{ouchi_uuid}

						// 交換されたご褒美を消化
						rewards.PUT("/digestion/:rewardExchangeId", controller.RewardDigestionHandler) // /v1/auth/users/ouchies/refresh/{ouchi_uuid}
					}

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

	// マルチパートフォームのメモリ使用制限を設定
	engine.MaxMultipartMemory = 8 << 20 // 20bit左シフトで8MiB

	// ルーティング
	routing(engine)

	// 静的ファイル設定
	loadingStaticFile(engine)

	// router設定されたengineを返す。
	return engine, nil
}
