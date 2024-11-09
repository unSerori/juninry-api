package view

import (
	"embed"
	"html/template"
	"io/fs"
	"juninry-api/common/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// テンプレートや静的ファイルを埋め込む
var (
	//go:embed views/*.html
	viewFilesFS embed.FS

	//go:embed views/styles/*.css
	stylesFS embed.FS

	//go:embed views/scripts/*.js
	scriptsFS embed.FS
)

// ファイルを設定
func LoadingStaticFile(engine *gin.Engine) error {
	// プロジェクト内のリソースをパスで指定するのではなく、埋め込んだembed.FS型変数で指定

	// views以下の静的ファイルを読み込む
	template, err := template.New("templateEmbedHTML"). // 新しいテンプレートを作成
								ParseFS( // 第一引数のFSから第二引数に該当するものを探し、テンプレートとする
			viewFilesFS,
			"views/*.html",
		) // エラーハンドル不要な場合は、template.Must()でラップすればパニックを発生させてくれる
	if err != nil {
		return err
	}
	engine.SetHTMLTemplate(template) // テンプレートを読み込む

	// 静的ファイルの調整
	adjustedStylesFS, err := fs.Sub(stylesFS, "views/styles") // 第一引数に対するサブファイルを、第二引数をルートとして作成
	if err != nil {
		logging.ErrorLog("Failed to create sub file.", err)
		panic(err)
	}
	adjustedScriptsFS, err := fs.Sub(scriptsFS, "views/scripts")
	if err != nil {
		logging.ErrorLog("Failed to create sub file.", err)
		panic(err)
	}

	// 静的ファイルを指定したURLで公開提供する(クライアントがアクセスするURL, サーバ上のリソース)
	engine.StaticFS("/styles", http.FS(adjustedStylesFS))
	engine.StaticFS("/scripts", http.FS(adjustedScriptsFS))

	logging.SuccessLog("Routing completed, start the server.")
	return nil
}
