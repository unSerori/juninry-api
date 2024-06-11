package service

import (
	"fmt"
	"time"
	"juninry-api/model"
)

type NoticeService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。



// userUuidからお知らせ一覧を持って来る
func (s *NoticeService) FindAllNotices(userUuid string) ([]model.Notice, error) {

	// userUuidを条件にしてclassUuidを取ってくる
	// 1 - userUuidからclass_menbershipの構造体を取ってくる
	classMemberships, err := model.FindClassMemberships(userUuid)
	if err != nil {
		return nil, err
	}

	// 2 - 構造体からclassUuidのスライス(配列)を作る
	var classUuids []string
	for _, classMembership := range classMemberships {
		classUuid := classMembership.ClassUuid
		classUuids = append(classUuids, classUuid)
	}

	// classUuidを条件にしてnoticeの構造体を取ってくる
	notices, err := model.FindNotices(classUuids)
	if err != nil {
		return nil, err
	}

	return notices, nil
}



// おしらせテーブル
type NoticeDetail struct { // typeで型の定義, structは構造体
	NoticeTitle       string	//お知らせのタイトル
	NoticeExplanatory string	//お知らせの内容
	NoticeDate        time.Time	//お知らせの作成日時
	UserName          string	// おしらせ発行ユーザ
	ClassName         string	// どのクラスのお知らせか
}

func (s *NoticeService) GetNoticeDetail(noticeUuid string) (model.Notice, error) {

	//お知らせ詳細情報取得
	noticeDetail, err := model.GetNoticeDetail(noticeUuid)
	if  err != nil {
		return model.Notice{}, err
	}

	fmt.Println(noticeDetail)


	//取ってきたnoticeDetailを整形します
	formattedNotice := NoticeDetail {
		NoticeTitle: noticeDetail.NoticeTitle,		//お知らせタイトル
		NoticeExplanatory: noticeDetail.NoticeExplanatory,	//お知らせの内容
		NoticeDate: noticeDetail.NoticeDate,			//お知らせ作成日時
	}


	return formattedNotice, err
}
