package batch

import "juninry-api/model"

func DeleteExpiredInviteCodes() {

	// 期限切れを削除する関数呼びます
	model.DeleteExpiredInviteCodes()
	

}
