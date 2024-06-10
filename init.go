package main

import (
	"fmt"
	"juninry-api/auth"
	"juninry-api/logging"
	"juninry-api/model"

	"github.com/joho/godotenv"
)

func Init() error {
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
		logging.ErrorLog("Error loading .env file", err)
		return err
	}

	// DB初期化
	err = model.DBConnect() // 接続
	if err != nil {
		logging.ErrorLog("Failed DB connect.", err)
		return err
	}
	err = model.MigrationTable() // テーブル作成
	if err != nil {
		logging.ErrorLog("Failed migration.", err)
		return err
	}
	// 接続を取得
	db := model.DBInstance()
	db.ShowSQL(true)       // SQL文の表示
	db.SetMaxOpenConns(10) // 接続数

	// 認証関連の初期化
	err = auth.InitAuth()
	if err != nil {
		logging.ErrorLog("Failed auth init.", err)
		return err
	}

	return nil
}
