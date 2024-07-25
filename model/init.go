package model

import (
	"fmt"
	"juninry-api/logging"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var db *xorm.Engine // インスタンス

// マイグレーション関連
func MigrationTable() error {
	// テーブルがないなら自動で作成。 // xormがテーブル作成時に列名をスネークケースにしてくれる。  // 列情報の追加変更は反映するが列の削除は反映しない。
	exists, _ := db.IsTableExist(&User{}) // この判定で、外部キー設定済みのテーブルの再Sync2時に外部キーのインデックスを消せないエラーを防いでいる。
	if !exists {
		err := db.Sync2(
			new(User),
			new(Ouchi),
			new(UserType),
			new(Class),
			new(ClassMembership),
			new(Notice),
			new(NoticeReadStatus),
			new(Subject),
			new(TeachingMaterial),
			new(Homework),
			new(HomeworkSubmission),
		)
		if err != nil {
			logging.ErrorLog("Failed to sync database.", err)
			return err
		}
	}

	// FK
	err := initFK()
	if err != nil {
		logging.ErrorLog("Failed to set foreign key.", err)
		return err
	}

	// サンプルデータ作成
	RegisterSample()

	return nil
}

// 外部キーを設定
func initFK() error {
	// User
	err := InitUserFK()
	if err != nil {
		return err
	}
	// ClassMembership
	err = InitClassMembershipFK()
	if err != nil {
		return err
	}
	// Notice
	err = InitNoticeFK()
	if err != nil {
		return err
	}
	// NoticeReadStatus
	err = InitNoticeReadStatus()
	if err != nil {
		return err
	}
	// TeachingMaterial
	err = InitTeachingMaterialFK()
	if err != nil {
		return err
	}
	// Homework
	err = InitHomeworkFK()
	if err != nil {
		return err
	}
	// HomeworkSubmission
	err = InitHomeworkSubmissionFK()
	if err != nil {
		return err
	}

	return err
}

// サンプルデータ作成
// 外部キーの参照先テーブルを先に登録する必要がある。
func RegisterSample() {
	// サンプル用データ作成
	CreateUserTypeTestData()
	CreateOuchiTestData()
	CreateClassTestData()
	CreateSubjectTestData()
	CreateTeachingMaterialTestData()
	// テスト用データ作成
	CreateUserTestData()
	CreateClassMembershipsTestData()
	CreateNoticeTestData()
	CreateNoticeReadStatusTestData()
	CreateHomeworkTestData()
	CreateHomeworkSubmissionTestData()
}

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

	return nil
}

// ORM初期化
func InitDB() (*xorm.Engine, error) {
	err := DBConnect() // DBインスタンスの生成とDBサーバー接続
	if err != nil {
		logging.ErrorLog("Failed DB connect.", err)
		return nil, err
	}
	err = MigrationTable() // テーブル作成
	if err != nil {
		logging.ErrorLog("Failed migration.", err)
		return nil, err
	}

	// 接続を取得
	db = DBInstance()

	// 設定
	db.ShowSQL(true)       // SQL文の表示
	db.SetMaxOpenConns(10) // 接続数

	fmt.Printf("db: %v\n", db)

	return db, nil
}

// 接続を取得
func DBInstance() *xorm.Engine {
	return db // 接続を返す
}
