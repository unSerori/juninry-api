// サービス内で発生したerrに名前を付けてcontroller側でswitchを使ったエラーハンドルをする
// 独自のエラー型構造体にはmsgとエラー型の情報を含む。エラー型情報も独自のタイプで、int管理のENUM
// サービス内でコントローラでswitch分岐させたいエラーが出たときはNewErrに紐づけたいエラー名とerr.Error()(:エラーmsg)を渡し、カスタムエラーを返す

package commons

// カスタムエラー型  // エラーの種類を示すErrTypeとエラーのmsgを持つ
type CustomErr struct {
	Type    ErrType
	Message string
}

// カスタムエラーのmsgを参照
func (e *CustomErr) Error() string {
	return e.Message
}

// ENUMでエラーの種類をまとめる
type ErrType int

const ( // ========================ここに新しい独自のエラーを追加していく
	ErrTypeHashingPassFailed ErrType = iota
	ErrTypeGenTokenFailed
	ErrTypeMITUKARANI
)

// エラー生成関数
func NewErr(errType ErrType, msg string) *CustomErr {
	return &CustomErr{
		Type:    errType,
		Message: msg,
	}
}
