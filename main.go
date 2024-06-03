package main // package

import (
	"fmt" // import
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/route"
)

// main method
func main() {
	err := Init() // 初期化処理
	if err != nil {
		fmt.Println("")
	}
	// 破棄処理
	defer logging.LogFile().Close()  // 関数終了時に破棄
	defer model.DBInstance().Close() // defer文でこの関数が終了した際に破棄する

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
