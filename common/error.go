// サービス内で発生したerrに名前を付けてcontroller側でswitchを使ったエラーハンドルをする
// 独自のエラー型構造体にはmsgとエラー型の情報を含む。エラー型情報も独自のタイプで、int管理のENUM
// サービス内で、コントローラでswitch分岐させたいエラーが出たときはNewErrに紐づけたいエラー名とerr.Error()(:エラーmsg)を渡し、カスタムエラーを返す

package common

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

// エラーの種類を定義
const ( // ========================ここに新しい独自のエラーを追加していく
	ErrTypeCustom            ErrType = iota
	ErrTypeHashingPassFailed         // ハッシュ化失敗
	ErrTypeGenTokenFailed            // トークン作成失敗
	ErrTypeNoResourceExist           // リソースが存在しない
	ErrTypePassMismatch              // パスワードが一致しない
)

// エラーに対するデフォルトmsgを設定
var errTypeMsg = map[ErrType]string{
	ErrTypeHashingPassFailed: "",
	ErrTypeGenTokenFailed:    "",
	ErrTypeNoResourceExist:   "could not find the relevant ID",
	ErrTypePassMismatch:      "password does not match",
}

// デフォルト引数をFunctional Option Patternで実装してみる

// 可変長引数で渡される省略可能なデフォルト引数たちを管理する構造体の定義
type NewErrParams struct {
	msg string
}

// オプション関数たちの返り血を定義、統一することで可変長引数としてまとめてforで処理できる
type NewErrParam func(*NewErrParams)

// デフォルト引数たちのデフォルト値を設定
// 設定後のデフォルト引数の構造体ポインタを返す
func defaultNewErrParams(errType ErrType) *NewErrParams {
	return &NewErrParams{
		msg: errTypeMsg[errType], // デフォルトの値を(mapから)設定
	}
}

// デフォルト引数msgのオプション関数、オプション関数はデフォルト引数を仕様とする関数の呼び出し側で使うのでパブリック
func WithMsg(msg string) NewErrParam {
	return func(nep *NewErrParams) { // デフォルト引数管理構造体を受け取り、その構造体にオプション関数が受け取った値を設定する無名関数を返す
		nep.msg = msg
	}
}

// エラー生成関数
// 第一引数にエラーの型、第二引数は設定するエラーメッセージを受け取る可変長引数（:可変長なので省略可能0~*）
func NewErr(errType ErrType, params ...NewErrParam) *CustomErr {
	// デフォルトを設定
	nep := defaultNewErrParams(errType)

	// 引数で渡された値があるなら代入
	for _, p := range params { // ？？？？関数型の可変長引数(:スライス)にparams(:可変長引数で渡される可能性がある引数をまとめた構造体)
		p(nep) // オプション関数が返した無名関数にデフォルト設定済みのデフォルト引数構造体を渡し、指定された値で上書きする=>引数にオプション関数と値をセットしないとデフォ設定済みの値が上書きされずに残る。
	}

	// 返す構造体ポインタを引数で受け取った値たちで初期化し作成、リターンする
	return &CustomErr{
		Type:    errType,
		Message: nep.msg,
	}
}
