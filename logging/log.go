package logging

import (
	"log"
	"time"
)

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
