package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logFile *os.File // ログファイル

// ログ初期設定側後にリソースを開放するために実態を返す
func LogFile() *os.File {
	return logFile // ログファイルを返す
}

// ログファイル出力のセットアップ
func SetupLogging() error {
	// ログファイルを作成
	var err error
	logFile, err = os.OpenFile("./logging/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil { // エラーチェック
		return fmt.Errorf("error opening file: %v", err) // エラーの場合
	}

	// ログの出力先をファイルにも。
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf("Set up logging.\n\n")

	return nil // ファイルを返す
}
