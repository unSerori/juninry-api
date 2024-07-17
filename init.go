package main

import (
	"fmt"
	"juninry-api/auth"
	"juninry-api/di"
	"juninry-api/logging"
	"juninry-api/model"
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

	// DB初期化
	err = model.DBConnect() // 接続
	if err != nil {
		logging.ErrorLog("Failed DB connect.", err)
		return nil, err
	}
	err = model.MigrationTable() // テーブル作成
	if err != nil {
		logging.ErrorLog("Failed migration.", err)
		return nil, err
	}
	// 接続を取得
	db := model.DBInstance()
	db.ShowSQL(true)       // SQL文の表示
	db.SetMaxOpenConns(10) // 接続数

	// 認証関連の初期化
	err = auth.InitAuth()
	if err != nil {
		logging.ErrorLog("Failed auth init.", err)
		return nil, err
	}

	// DIコンテナ関連
	initInstances.Container = di.BuildContainer()

	// スケジューラを初期化して開始
	scheduler.StartScheduler()

	return initInstances, nil
}
