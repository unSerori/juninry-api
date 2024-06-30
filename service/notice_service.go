package service

import (
	"fmt"
	"juninry-api/auth"
	"juninry-api/logging"
	"juninry-api/model"
	"juninry-api/common"
	"time"

	"github.com/google/uuid"
)

type NoticeService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// noticeの新規登録
func (s *NoticeService) RegisterNotice(bNotice model.Notice) (string, error) {
	// notice_uuidを生成
	noticeId, err := uuid.NewRandom() //新しいuuidの作成
	if err != nil {
		return "", err
	}
	bNotice.NoticeUuid = noticeId.String() //設定

	// 構造体をレコード登録処理に投げる
	_, err = model.CreateNotice(bNotice) // 第一返り血は登録成功したレコード数
	if err != nil {
		return "", err
	}

	// 登録が成功したらトークンを生成する
	token, err := auth.GenerateToken(bNotice.NoticeUuid) // トークンを取得
	if err != nil {
		return "", common.NewErr(common.ErrTypeGenTokenFailed, common.WithMsg(err.Error()))
	}

	return token, nil
}

// おしらせテーブル(1件取得用)
type NoticeDetail struct { // typeで型の定義, structは構造体
	NoticeTitle       string    //お知らせのタイトル
	NoticeExplanatory string    //お知らせの内容
	NoticeDate        time.Time //お知らせの作成日時
	UserName          string    // おしらせ発行ユーザ
	ClassName         string    // どのクラスのお知らせか
}

// お知らせ詳細取得
func (s *NoticeService) GetNoticeDetail(noticeUuid string) (NoticeDetail, error) {

	//お知らせ詳細情報取得
	noticeDetail, err := model.GetNoticeDetail(noticeUuid)
	if err != nil {
		return NoticeDetail{}, err //nilで返せない!不思議!!
	}

	//確認用
	// fmt.Println(noticeDetail)

	//取ってきたnoticeDetailを整形して、controllerに返すformatに追加する
	formattedNotice := NoticeDetail{
		NoticeTitle:       noticeDetail.NoticeTitle,       //お知らせタイトル
		NoticeExplanatory: noticeDetail.NoticeExplanatory, //お知らせの内容
		NoticeDate:        noticeDetail.NoticeDate,        //お知らせ作成日時
	}

	//userUuidをuserNameに整形
	userUuid := noticeDetail.UserUuid
	user, nil := model.GetUser(userUuid)
	if err != nil {
		return NoticeDetail{}, err
	}
	//整形後formatに追加
	formattedNotice.UserName = user.UserName // おしらせ発行ユーザ

	//classUuidをclassNameに整形
	classUuid := noticeDetail.ClassUuid
	class, nil := model.GetClass(classUuid) //←構造体で帰ってくる！！
	if err != nil {
		return NoticeDetail{}, err
	}
	//整形後formatに追加
	formattedNotice.ClassName = class.ClassName // どのクラスのお知らせか

	return formattedNotice, err
}

// おしらせテーブル(全件取得用)
type Notice struct { // typeで型の定義, structは構造体
	NoticeTitle string    //お知らせのタイトル
	NoticeDate  time.Time //お知らせの作成日時
	UserName    string    // おしらせ発行ユーザ
	ClassName   string    // どのクラスのお知らせか
	ReadStatus  int       //お知らせを確認しているか
}

// ユーザの所属するクラスのお知らせ全件取得
func (s *NoticeService) FindAllNotices(userUuid string) ([]Notice, error) {

	// userUuidを条件にしてclassUuidを取ってくる
	// 1 - userUuidからclass_membershipの構造体を取ってくる
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

	// TODO:データを逆順に追加するために一時的なスライス
	var temp []Notice

	//noticesの一つをnoticeに格納(for文なのでデータ分繰り返す)
	for _, notice := range notices {

		//整形する段階で渡されるuserUuidが消えてしまうため、saveに保存しておく
		userUuidSave := userUuid

		//noticeを整形して、controllerに返すformatに追加する
		notices := Notice{
			NoticeTitle: notice.NoticeTitle, //お知らせのタイトル
			NoticeDate:  notice.NoticeDate,  //お知らせの作成日時
		}

		//userUuidをuserNameに整形
		userUuid := notice.UserUuid
		user, nil := model.GetUser(userUuid) //ユーザ取得
		if err != nil {
			return []Notice{}, err
		}

		//整形後formatに追加
		notices.UserName = user.UserName // おしらせ発行ユーザ

		//classUuidをclassNameに整形
		classUuid := notice.ClassUuid
		class, nil := model.GetClass(classUuid) //クラス取得
		if err != nil {
			return []Notice{}, err
		}
		//整形後formatに追加
		notices.ClassName = class.ClassName // おしらせ発行ユーザ

		//確認しているか取得
		status, err := model.GetNoticeReadStatus(notice.NoticeUuid, userUuidSave)
		if err != nil {
			return []Notice{}, err
		}

		fmt.Println(status)
		//確認していた場合、ReadStatusに1を保存する
		if status {
			notices.ReadStatus = 1
		} else {
			notices.ReadStatus = 0
		}

		//宣言したスライスに追加していく
		// formattedAllNotices = append(formattedAllNotices, notices)
		temp = append(temp, notices) //並べ替えるために一時的にtempに保存する
	}

	//fomat後のnotices格納用変数(複数返ってくるのでスライス)
	var formattedAllNotices []Notice

	// tempを逆順にしてformattedAllNoticesに追加する
	for i := len(temp) - 1; i >= 0; i-- {
		formattedAllNotices = append(formattedAllNotices, temp[i])
	}

	//確認用
	fmt.Println(formattedAllNotices)

	return formattedAllNotices, nil
}

// noticeの既読登録
func (s *NoticeService) ReadNotice(bRead model.NoticeReadStatus) error {

  // クラス作成権限を持っているか確認
	isTeacher, err := model.IsParent(bRead.UserUuid)
	if err != nil { // エラーハンドル
		return err
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return common.NewErr(common.ErrTypePermissionDenied)
	}

	// 構造体をレコード登録処理に投げる
	err = model.ReadNotice(bRead) // 第一返り血は登録成功したレコード数
	if err != nil {
		return err
	}

	return nil
}
