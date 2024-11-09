package logging

import (
	"fmt"
	"log"
	"time"
)

// 成功時のログをログファイルに残す
// successName: 先頭大文字コロンあり
func SuccessLog(successName string) {
	log.Printf("Success LOG: %s\n", successName)
	log.Printf("Time: %v\n\n", time.Now()) // 時刻
}

// エラー時のログをログファイルに残す
// errName: 先頭大文字コロンあり
// err: err or 先頭小文字コロンなし
func ErrorLog(errName string, err error) {
	log.Printf("ERROR LOG: %s\n", errName)
	log.Printf("Time: %v\n", time.Now()) // 時刻
	if err != nil {
		log.Printf("Error: %s\n\n", err)
	} else {
		log.Printf("Error: NIL")
	}
}

// 情報の記録
// title: 先頭大文字コロンあり
// info: 空文字 or 先頭大文字コロンあり
func InfoLog(title string, info string) {
	log.Printf("INFO LOG: %s\n", title)
	log.Printf("Time: %v\n", time.Now()) // 時刻
	if info != "" {
		log.Printf("Info: \n%s\n", info)
		//log.Printf("Info: %s\n", info)
	}
}

// 警告
// title: 先頭大文字コロンあり
// warning: 先頭大文字コロンあり
func WarningLog(title string, warning string) {
	log.Printf("WARNING LOG: %s\n", title)
	log.Printf("Time: %v\n", time.Now()) // 時刻
	if warning != "" {
		log.Printf("Warning: %s\n\n", warning)
	}
}

// 単純なprintf
// args: 結合される
func SimpleLog(args ...interface{}) {
	// 結合後の変数
	var message string

	// forで引数を接続
	for _, arg := range args {
		message += fmt.Sprintf("%v", arg)
	}

	// ログに書き込み
	log.Print(message)
}
