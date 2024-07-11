package batch

import "juninry-api/model"

func DeleteExpiredInviteCodes() {

	// 期限切れを削除する関数呼びます(クラス招待コード)
	model.DeleteExpiredInviteCodes()
	
	// 期限切れを削除する関数呼びます２(おうち招待コード)
	model.DeleteExpiredOuchiInviteCodes()

}
