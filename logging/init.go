package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logFile *os.File // ログファイル

// ログファイルを作成
func openLogFile() error {
	var err error
	logFile, err = os.OpenFile("./logging/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil { // エラーチェック
		return fmt.Errorf("error opening file: %v", err) // エラーの場合
	}
	return nil
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
