package main // package

import (
	"fmt" // import
	"juninry-api/route"
)

// main method
func main() {
	// router設定されたengineを受け取る。
	router, err := route.GetRouter()
	if err != nil {
		fmt.Println(err) // エラー内容を出力し早期リターン ログ関連のエラーなのでログは出力しない
		return
	}

	// テンプレートと静的ファイルを読み込む
	router.LoadHTMLGlob("view/views/*.html")
	router.Static("/styles", "./view/views/styles") // クライアントがアクセスするURL, サーバ上のパス
	router.Static("/scripts", "./view/views/scripts")

	// 鯖起動
	router.Run(":4561")
}
