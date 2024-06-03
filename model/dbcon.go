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
func DBConnect() error {
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
	exists, _ := db.IsTableExist(&User{}) // この判定で、外部キー設定済みのテーブルの再Sync2時に外部キーのインデックスを消せないエラーを防いでいる。
	if !exists {
		err = db.Sync2(
			new(User),
			new(Ouchi),
			new(UserType),
		)
		if err != nil {
			log.Fatalf("Failed to sync database: %v", err)
			return err
		}

	}

	// FK
	err = initFK()
	if err != nil {
		fmt.Println(err)
		fmt.Println("NNNNNNNNNNNNNNNNNNNNNN")
		return err
	}

	// テスト用データ作成
	CreateUserTypeTestData()
	CreateOuchiTestData()

	return nil
}

// 接続を取得
func DBInstance() *xorm.Engine {
	return db // 接続を返す
}

// 外部キーを設定
func initFK() error {
	// User
	err := InitUserFK()
	if err != nil {
		return err
	}

	return err
}
