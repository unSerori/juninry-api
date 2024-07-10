package scheduler

import (
	"juninry-api/batch"
	"juninry-api/logging"

	"github.com/robfig/cron/v3"
)

// スケジューラーを開始
func StartScheduler() {
	cron := cron.New()

	// 処理する関数と時間を設定

	// 4時に期限切れ招待コードを破壊する
	_, err := cron.AddFunc("0 4 * * *", batch.DeleteExpiredInviteCodes)
	if err != nil { //エラーハンドル
		logging.ErrorLog("Class creation was not possible due to other problems.", err)
		panic(err)
	}

	cron.Start()
}
