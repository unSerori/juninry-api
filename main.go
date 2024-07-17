package main // package

import (
	// import

	"fmt"
	"juninry-api/logging"
	"juninry-api/model"

	"github.com/gin-gonic/gin"
)

// main method
func main() {
	// 初期化処理
	initInstances, err := Init() // add "initInstances, " when changing to ddd
	if err != nil {
		return
	}
	// 破棄処理
	defer logging.LogFile().Close()  // 関数終了時に破棄
	defer model.DBInstance().Close() // defer文でこの関数が終了した際に破棄する
	logging.SuccessLog("Successful server init process.")

	// // router設定されたengineを受け取る。
	// router, err := route.SetupRouter()
	// if err != nil {
	// 	logging.ErrorLog("Couldnt receive router engine.", err) // エラー内容を出力し早期リターン
	// 	return
	// }

	// 鯖起動  // router.Run(":4561")
	err = initInstances.Container.Invoke( // 依存性注入コンテナから必要な依存解決を解決し、渡されたコールバック関数にcontainerが持つエンジンの実体を渡す
		func(r *gin.Engine) { // containerが持つエンジンを受け取り鯖を起動
			r.Run(":4561") // 指定したポートで鯖起動
		},
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
}
