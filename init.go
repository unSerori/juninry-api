package main

import (
	"fmt"
	"juninry-api/logging"
	"juninry-api/model"
	"log"

	"github.com/joho/godotenv"
)

func Init() error {
	// .envから定数をプロセスの環境変数にロード
	err := godotenv.Load(".env") // エラーを格納
	if err != nil {              // エラーがあったら
		//logging.ErrorLog("Error loading .env file", err)
		panic("Error loading .env file.")
	}

	// ログ設定を初期化
	err = logging.SetupLogging() // セットアップ
	if err != nil {              // エラーチェック
		fmt.Printf("error opening file: %v\n", err)
	}
	log.Println("Start server!")

	// DB初期化
	err = model.DBConnect() // 接続
	if err != nil {
		fmt.Println("")
	}
	// 接続を取得
	db := model.DBInstance()
	db.ShowSQL(true)       // SQL文の表示
	db.SetMaxOpenConns(10) // 接続数

	return nil
}
