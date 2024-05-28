package model

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var db *xorm.Engine // インスタンス

// SQL接続とテーブル作成
func DBConnect() {
	// 環境変数から取得
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbDB := os.Getenv("MYSQL_DATABASE")

	// Mysqlに接続
	var err error
	db, err = xorm.NewEngine( // dbとエラーを取得
		"mysql", // dbの種類"root:root@tcp(db:3306)/cgroup?charset=utf8"
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbDB), // 接続情報
	)
	if err != nil { // エラー処理
		fmt.Println("せつぞくできなかった")
		log.Fatal("Couldnt connect to the db server.", err)
	} else {
		fmt.Println("せつぞくできた")
		log.Println("Could connect to the db server.")
	}

	// テーブルがないなら自動で作成。 // xormがテーブル作成時にモデル名を複数形に、列名をスネークケースにしてくれる。  // 列情報の追加変更は反映するが列の削除は反映しない。
	err = db.Sync2(
		new(User),
	)
	if err != nil {
		log.Fatalf("Failed to sync database: %v", err)
	}

	// テスト用データ作成
	// CreateKHogeTestData()
}

// 接続を取得
func DBInstance() *xorm.Engine {
	return db // 接続を返す
}
