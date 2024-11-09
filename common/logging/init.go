package logging

import (
	"fmt"
	"io"
	"juninry-api/utility"
	"log"
	"os"
)

var logFile *os.File // ログファイル

// ログファイルを作成
func openLogFile() error {
	// ディレクトリを作成
	err := utility.SafeMkdir("./common/logging", 0755, ErrorLog)
	if err != nil {
		return err
	}

	// ログファイルの作成
	logFile, err = os.OpenFile("./common/logging/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil { // エラーチェック
		return fmt.Errorf("error opening file: %v", err) // エラーの場合
	}
	return nil
}

// ログ初期設定側後にリソースを開放するために実態を返す
func LogFile() *os.File {
	return logFile // ログファイルを返す
}

// 出力先変更
func SetupLogOutput() {
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf("Set up logging.\n\n")
}

// ログファイル出力のセットアップ
func InitLogging() error {
	// ログファイルを作成
	err := openLogFile()
	if err != nil {
		return err
	}

	// ログの出力先をファイルにも。
	SetupLogOutput()

	return nil // ファイルを返す
}
