// リポジトリインターフェース

package domain

import (
	"mime/multipart"
	"os"
)

// リポジトリインターフェース  // 永続化層ですべて実装する。永続化層が一つのパッケージにまとめられていない場合、それぞれのパッケージに対してのインターフェースを提供しておく。

// ファイルディレクトリ操作
type TeachingMaterialRepository interface {
	// ファイル操作
	CreateDstDir(dst string, fileMode os.FileMode) error    // 指定されたディレクトリを作成
	UpLoadImage(filePath string, file multipart.File) error // ファイルをディレクトリに保存

	// DB操作
	GetPermissionInfoById(id string) (int, error)              // idから権限情報を取得
	IsUserInClass(classId string, userId string) (bool, error) // ユーザーがクラスに属しているかどうか
	IsSubjectExists(subjectId int) (bool, error)               // それが登録済みの教科かどうか
	CreateTM(tm TeachingMaterial) error                        // 完成した教材構造体をインサート
}
