package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var logFile *os.File // ログファイル

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

// ログ初期設定側後にリソースを開放するために実態を返す
func LogFile() *os.File {
	return logFile // ログファイルを返す
}

// 成功時のログをログファイルに残す
func SuccessLog(successName string) {
	log.Printf("Success LOG: %s\n", successName)
	log.Printf("Time: %v\n\n", time.Now()) // 時刻
}

// エラー時のログをログファイルに残す
func ErrorLog(errName string, err error) {
	log.Printf("ERROR LOG: %s\n", errName)
	log.Printf("Time: %v\n", time.Now()) // 時刻
	log.Printf("Error: %s\n\n", err)
}
