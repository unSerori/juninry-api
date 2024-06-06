package controller

import (
	"fmt"
	"juninry-api/logging"
	"juninry-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var homeworkService = service.HomeworkService{}

// 課題全件取得
func FindHomeworkHandler(c *gin.Context) {
	// ユーザーを特定する
	// id, exists := c.Get("id")
	// if !exists { // idがcに保存されていない。
	// 	// エラーログ
	// 	logging.ErrorLog("The id is not stored.", nil)
	// 	// レスポンス
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"srvResCode": 7013,                    // コード
	// 		"srvResMsg":  "The id is not stored.", // メッセージ
	// 		"srvResData": gin.H{},                 // データ
	// 	})
	// 	return
	// }

	//テスト用userUuid割り当て
	id := "3cac1684-c1e0-47ae-92fd-6d7959759224"

	//TODO: エラーのハンドリングがカス
	//問い合わせ処理と失敗レスポンス
	homeworks, err := homeworkService.FindHomework(id)
	if err != nil { //エラーハンドル
		//DB関連のエラー
		return
	}


	//構造体の配列をレスポンスの形に整形
	var homeworkList []gin.H
	for _, homework := range homeworks {
		homeworkJSON := gin.H{
			"taskLimit":          homework.HomeworkLimit,      	//タスクの期限
			"startPage":          homework.StartPage,          	// 開始ページ
			"PageCount":          homework.PageCount,          	// ページ数
			"homeworkPosterUUID": homework.HomeworkPosterUuid, 	// 投稿者ID
			"homeworkNote":       homework.HomeworkNote,       	// 宿題説明
		}
		homeworkList = append(homeworkList, homeworkJSON)		//リストにぶちこむ
	}

	fmt.Println(homeworkList)
	// 処理後の成功
	// 成功ログ
	logging.SuccessLog("Successful get homework.")
	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		"srvResCode": 1001,
		"srvResMsg":  "Successful get homework.",
		"srvResData": homeworkList,
	})

}
