package requests

// ハードデバイスの初期化
// 必要な項目が増えたら増やす
// 送らないデバイスからは暗黙的にnullが送られるはず
type InitHard struct {
	HardwareTypeId int    `json:"hardwareTypeId"`
	OuchiUuid      string `json:"ouchiUUID"`
}
