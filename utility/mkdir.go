// その他のutility

package utility

import (
	"os"
)

// MkdirAllを使ってディレクトリを再起的に作成する
func SafeMkdir(dirPath string, fileMode os.FileMode, logError func(string, error)) error { // fileMode: 0644
	// ディレクトリの存在チェックを行い、存在しないなら作成
	if _, err := os.Stat(dirPath); os.IsNotExist(err) { // ファイル情報を取得し、失敗するとerrが返る // errの原因がファイルが存在しないことなら新しく作成
		if err := os.MkdirAll(dirPath, fileMode); err != nil {
			logError("Create dir: Failed to create.", err)
		}
	} else if err != nil { // ディレクトリが存在するにも関わらずエラーが発生した時
		// logging.ErrorLog("Create dir: Errors other than IsNotExist.", err)
		logError("Create dir: Errors other than IsNotExist.", err)
		return err
	}

	return nil
}
