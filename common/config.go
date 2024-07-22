// 環境変数からロードしたいとき。
package common

import (
	"fmt"
	"juninry-api/logging"
	"os"
	"strconv"
)

func LoadReqBodyMaxSize(defaultSize int64) int64 {
	maxSize := defaultSize // デフォ値を設定、.envの環境変数がなければこれがそのまま返る
	fmt.Printf("maxSize: %v\n", maxSize)
	if maxSizeByEnv := os.Getenv("REQ_BODY_MAX_SIZE"); maxSizeByEnv != "" { // 空文字でなければ数値に変換する
		maxSizeByEnvInt, err := strconv.Atoi(maxSizeByEnv) // 数値に変換
		fmt.Printf("maxSizeByEnvInt: %v\n", maxSizeByEnvInt)
		if err != nil {
			fmt.Println("aa")
			logging.ErrorLog("Numerical conversion of environment variables in LoadReqBodyMaxSize(defaultSize int64) fails.", err)
		} else {
			fmt.Println("bb")
			maxSize = int64(maxSizeByEnvInt)
		}
	}
	fmt.Printf("maxSize: %v\n", maxSize)
	return maxSize
}
