package main

import (
	"fmt"
	"juninry-api/auth"
	"juninry-api/logging"
	"juninry-api/route"
	"juninry-api/scheduler"

	"go.uber.org/dig"

	"github.com/joho/godotenv"
)

// 初期化の成果物
type InitInstances struct {
	Container *dig.Container
}

// mainでの初期化処理
func Init() (*InitInstances, error) {
	// 結果
	// var initInstances *InitInstances  // ポインタ型の宣言(不要)
	initInstances := &InitInstances{} // 同じ: initInstances := new(InitInstances)

	// ログ設定を初期化
	err := logging.InitLogging() // セットアップ
	if err != nil {              // エラーチェック
		fmt.Printf("error set up logging: %v\n", err) // ログ関連のエラーなのでログは出力しない
		panic("error set up logging.")
	}
	logging.SuccessLog("Start server!")

	// .envから定数をプロセスの環境変数にロード
	err = godotenv.Load(".env") // エラーを格納
	if err != nil {             // エラーがあったら
		logging.ErrorLog("Error loading .env file.", err)
		return nil, err
	}

	// 認証関連の初期化
	err = auth.InitAuth()
	if err != nil {
		logging.ErrorLog("Failed auth init.", err)
		return nil, err
	}

	// DB初期化やルーティング設定など、依存関係にかかわるものの初期化とDIコンテナによる各層の依存関係登録
	initInstances.Container = route.BuildContainer()

	// スケジューラを初期化して開始
	scheduler.StartScheduler()

	return initInstances, nil
}
