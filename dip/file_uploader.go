// ファイルアップローダ

package dip

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

// 依存性の逆転(DIP)のための、インターフェース定義
// 実際に使う関数を、具体的な実装(gin.ContextとそれをつかうController層)への依存から、抽象的なインターフェースに依存させるために必要。
// このインターフェースはDIPしたい機能を提供する。このインターフェースを通じて逆転させたい機能を利用する
type FileUpLoader interface {
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
}

// 依存性の逆転(DIP)のための、ラッパー構造体定義
// 実際の実装で元機能にアクセスするために必要。
// 具体的な実装をラップし、これを介して抽象的な実装を実現する
type GinContextWrapper struct {
	C *gin.Context
}

// 依存性の逆転(DIP)のための、実装された関数
// interfaceを関数に実装するために必要。
// ラッパー構造体を介して元の機能を呼び出す
func (gcw *GinContextWrapper) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	return gcw.C.SaveUploadedFile(file, dst) // 元関数を呼び出す
}
