package scheduler

import (
	"juninry-api/batch"
	"juninry-api/logging"

	"github.com/robfig/cron/v3"
)

// スケジューラーを開始
func StartScheduler() {
	cron := cron.New()

	// エラーが入る君
	var err error

	// 4時に期限切れクラス招待コードを破壊する
	_, err = cron.AddFunc("0 4 * * *", batch.DeleteExpiredInviteCodes)
	if err != nil { //エラーハンドル
		logging.ErrorLog("Class creation was not possible due to other problems.", err)
		panic(err)
	}

	// 4時に期限切れおうち招待コードを破壊
	_, err = cron.AddFunc("0 4 * * *", batch.DeleteExpiredOuchiInviteCodes)
	if err != nil { //エラーハンドル
		logging.ErrorLog("Ouchi creation was not possible due to other problems.", err)
		panic(err)
	}


	cron.Start()
}
