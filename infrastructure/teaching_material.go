// サービスから呼び出される具体的な処理関数

package infrastructure

import (
	"errors"
	"io"
	"juninry-api/domain"
	"juninry-api/model"
	"juninry-api/utility/custom"
	"mime/multipart"
	"os"

	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// 永続化の構造体  // 依存先はdomain/repositoryと永続化処理に必要なインスタンス。repositoryはinterfaceなのでここには書かず、interfaceを実際に実装する。
type TeachingMaterialPersistence struct {
	db *xorm.Engine
}

// ファクトリー関数
func NewTeachingMaterialPersistence(db *xorm.Engine) domain.TeachingMaterialRepository {
	return &TeachingMaterialPersistence{db: db} // 返す構造体インスタンスのメソッドは、ファクトリー関数の返り血のインターフェースをすべて実装している(implements)ので、型が違うが無問題
}

// このエンティティのリポジトリインターフェースをすべて実装

// ファイル操作

// 指定されたディレクトリを作成
func (p *TeachingMaterialPersistence) CreateDstDir(dst string, fileMode os.FileMode) error { // dst = "./hoge/piyo"
	// ディレクトリが存在しない場合
	if _, err := os.Stat(dst); os.IsNotExist(err) { // ファイル情報を取得, 取得できないならerrができる // 取得できなかったとき、ファイルが存在しないことが理由なら新しく作成
		if err := os.MkdirAll(dst, fileMode); err != nil { // fileMode: 0644
			return err
		}
	}
	return nil
}

// ファイルをディレクトリに保存
func (p *TeachingMaterialPersistence) UpLoadImage(filePath string, file multipart.File) error {
	oFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644) // ファイルが存在しない場合に新規作成|O_CREATEと組み合わせることで同名ファイル存在時にエラーを発生|書き込み専用で開く
	if err != nil {
		return err
	}
	defer oFile.Close() // リソース解放
	// 読み書き位置の設定
	if _, err := file.Seek(0, io.SeekStart); err != nil { // 書き込みたいデータ
		return err
	}
	if _, err := oFile.Seek(0, io.SeekStart); err != nil { // 開いたファイル
		return err
	}
	// データをコピー
	if _, err := io.Copy(oFile, file); err != nil { // io.Copy()はimage<-*multipart.FileHeaderを解釈できないので、バイナリからファイルタイプを特定するために取得したFileオブジェクトを利用
		return err
	}

	return nil
}

// DB操作

// ユーザーのuuidからユーザー権限を取得
func (p *TeachingMaterialPersistence) GetPermissionInfoById(id string) (int, bool, error) {
	var user model.User // 取得したデータをマッピングする構造体
	isFound, err := p.db.Where("user_uuid = ?", id).Get(&user)
	if err != nil || !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`err != nil || !isFound`で、影響を与えるSQL文の時は`err != nil || affected == 0`でハンドリング
		return 0, false, err
	}
	return user.UserTypeId, true, nil // teacher, junior, patron
}

// ユーザーがクラスに属しているかどうか
func (p *TeachingMaterialPersistence) IsUserInClass(classId string, userId string) (bool, error) {
	var cm model.ClassMembership // 取得したデータをマッピングする構造体
	isFound, err := p.db.Where("class_uuid = ? AND user_uuid = ?", classId, userId).Get(&cm)
	if err != nil || !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`err != nil || !isFound`で、影響を与えるSQL文の時は`err != nil || affected == 0`でハンドリング
		return false, err
	}
	return true, nil
}

// それが登録済みの教科かどうか
func (p *TeachingMaterialPersistence) IsSubjectExists(subjectId int) (bool, error) {
	var subject model.Subject // 取得したデータをマッピングする構造体
	isFound, err := p.db.Where("subject_id = ?", subjectId).Get(&subject)
	if err != nil || !isFound { //エラーハンドル  // 影響を与えないSQL文の時は`err != nil || !isFound`で、影響を与えるSQL文の時は`err != nil || affected == 0`でハンドリング
		return false, err
	}
	return true, nil
}

// 完成した教材構造体をインサート
func (p *TeachingMaterialPersistence) CreateTM(tm domain.TeachingMaterial) error {
	// エンティティをテーブルモデル構造体に変換
	tmModel := model.FromDomainEntity(&tm)

	// INSERT
	affected, err := p.db.Insert(tmModel)
	if err != nil { //エラーハンドル
		// XormのORMエラーを仕分ける
		var mysqlErr *mysql.MySQLError // DBエラーを判定するためのDBインスタンス
		if errors.As(err, &mysqlErr) { // errをmysqlErrにアサーション出来たらtrue
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反
				return custom.NewErr(custom.ErrTypeUniqueConstraintViolation)
			default: // ORMエラーの仕分けにぬけがある可能性がある
				return custom.NewErr(custom.ErrTypeOtherErrorsInTheORM)
			}
		}

		// 通常の処理エラー
		return err // 受け取ったエラーを返す
	}
	if affected == 0 {
		return custom.NewErr(custom.ErrTypeNoFoundR)
	}
	return nil
}
