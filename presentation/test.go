// ユーザーインターフェース(:リクエストの受け取りとレスポンスの返却)

package presentation

import (
	"bytes"
	"fmt"
	"io"
	"juninry-api/common/logging"
	"juninry-api/utility"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// /
func ShowRootPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"topTitle":  "Route /",                                                            // ヘッダ内容 h1
		"mainTitle": "Hello.",                                                             // メインのタイトル h2
		"time":      time.Now(),                                                           // 時刻
		"message":   "This is an API server written in Golang for safety check purposes.", // message
	})
}

// cfmreq
func ConfirmationReq(c *gin.Context) {
	// 値をセット
	cmfReqValues := struct { // jsonタグを使ってc.JSON()で構造体をJSONとして返却するときのためにkey値を指定する
		Method      string              `json:"method"`
		Url         *url.URL            `json:"url"`
		Header      http.Header         `json:"header"`
		PathPrams   map[string]string   `json:"pathPrams"`
		QueryParams map[string][]string `json:"queryParams"`
		Body        string              `json:"body"`
	}{
		Method: c.Request.Method,
		Url:    c.Request.URL,
		Header: c.Request.Header,
		PathPrams: func() map[string]string {
			pathParams := make(map[string]string)
			for _, param := range c.Params {
				pathParams[param.Key] = param.Value
			}
			return pathParams
		}(), // 無名関数は定義後間髪入れず()引数を渡して呼び出す
		QueryParams: c.Request.URL.Query(),
		Body: func() string {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logging.ErrorLog("Faild to reading req body.", err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			return string(body)
		}(),
	}

	// サーバーデバッグコンソールで確認
	utility.CheckStruct(cmfReqValues)
	fmt.Println()

	// 成功ログ
	logging.SuccessLog("JSON for testing.")
	// レスポンス
	resStatusCode := http.StatusOK
	c.JSON(http.StatusOK, gin.H{ // bodyがJSON形式のレスポンスを返す
		"srvResCode": http.StatusText(resStatusCode), // メッセージ
		"srvResData": gin.H{
			"message": "hello go server!",
			"info":    cmfReqValues,
		}, // データ
	})
}

// test
func Test(c *gin.Context) {
	// URLに含める時刻形式の吟味
	date := "2024-10-13T20C58C55P09C00"
	fmt.Printf("date: %v\n", date)
	dateAdjusted := strings.ToUpper(date)
	dateAdjusted = strings.Replace(dateAdjusted, "C", ":", -1) // :をURLで送りたくないため
	dateAdjusted = strings.Replace(dateAdjusted, "M", "-", -1) // UTCからの負方向の時差
	dateAdjusted = strings.Replace(dateAdjusted, "P", "+", -1) // UTCからの正方向の時差
	fmt.Printf("dateAdjusted: %v\n", dateAdjusted)
	dateTT, err := time.Parse(time.RFC3339, dateAdjusted) // RFC3339: "2006-01-02T15:04:05Z07:00"
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("dateTT: %v\n", dateTT)

	// 文字数
	passes := []string{
		"aaaaaaaaaaa",
		"aaaaaaaaaaaa",
		"aaaaaaaaaaaaa",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",

		"あああああああああああああああああああああああ",
		"ああああああああああああああああああああああああ",
		"あああああああああああああああああああああああああ",
	}

	for _, pass := range passes {
		fmt.Print(pass, ": ", len(pass), "\n")
		fmt.Print(pass, ": ", utf8.RuneCountInString(pass), "\n")
		fmt.Println()
	}
	// emails := []string{
	// 	"hoge@gmail.com",
	// 	"piyo.ta@gmail.com",
	// 	"piyo-ta@gamil.com",
	// 	"tyu320v9",
	// 	"8898@g.c",
	// 	"---@g.com",
	// 	"hoge@piyo",
	// 	"..@a",
	// 	"a@.",
	// }

	// for _, email := range emails {
	// 	_, err := mail.ParseAddress(email)
	// 	if err != nil {
	// 		fmt.Println(email + ": " + "no")
	// 	} else {
	// 		fmt.Println(email + ": " + "ok")
	// 	}
	// }
}
