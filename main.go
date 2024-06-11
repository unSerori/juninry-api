package main // package

import (
	// import
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/route"
)

// main method
func main() {
	// 初期化処理
	err := Init()
	if err != nil {
		return
	}
	// 破棄処理
	defer logging.LogFile().Close()  // 関数終了時に破棄
	defer model.DBInstance().Close() // defer文でこの関数が終了した際に破棄する
	logging.SuccessLog("Successful server init process.")

	// router設定されたengineを受け取る。
	router, err := route.SetupRouter()
	if err != nil {
		logging.ErrorLog("Couldnt receive router engine.", err) // エラー内容を出力し早期リターン
		return
	}

	// 鯖起動
	router.Run(":4561")
}
