package service

import (
	"fmt"
	"juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"time"

	"github.com/google/uuid"
)

type NoticeService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// noticeの新規登録
func (s *NoticeService) RegisterNotice(bNotice model.Notice) error {

	//先生かのタイプチェック
	isTeacher, err := model.IsTeacher(bNotice.UserUuid)
	if err != nil { // エラーハンドル
		return err
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return common.NewErr(common.ErrTypePermissionDenied)
	}

	// notice_uuidを生成
	noticeId, err := uuid.NewRandom() //新しいuuidの作成
	if err != nil {
		return err
	}
	bNotice.NoticeUuid = noticeId.String() //設定

	// 構造体をレコード登録処理に投げる
	_, err = model.CreateNotice(bNotice) // 第一返り血は登録成功したレコード数
	if err != nil {
		return err
	}

	return nil

}

// おしらせテーブル(1件取得用)
type NoticeDetail struct { // typeで型の定義, structは構造体
	NoticeUuid        string    // おしらせUUID
	NoticeTitle       string    // おしらせのタイトル
	NoticeExplanatory string    // おしらせの内容
	NoticeDate        time.Time // おしらせの作成日時
	UserName          string    // おしらせ発行ユーザ
	ClassUuid         string    // クラスUUID
	ClassName         string    // どのクラスのお知らせか
	RefUuid           string    // 引用UUID
	ReadStatus        int       // 既読フラグ
}

// お知らせ詳細取得
func (s *NoticeService) GetNoticeDetail(noticeUuid string, userUuid string) (NoticeDetail, error) {

	//お知らせ詳細情報取得
	noticeDetail, err := model.GetNoticeDetail(noticeUuid)
	if err != nil {
		return NoticeDetail{}, err //nilで返せない!不思議!!
	}

	//取ってきたnoticeDetailを整形して、controllerに返すformatに追加する
	formattedNotice := NoticeDetail{
		NoticeTitle:       noticeDetail.NoticeTitle,       //お知らせタイトル
		NoticeExplanatory: noticeDetail.NoticeExplanatory, //お知らせの内容
		NoticeDate:        noticeDetail.NoticeDate,        //お知らせ作成日時
		NoticeUuid:        noticeDetail.NoticeUuid,        //おしらせ引用UUID
		ClassUuid:         noticeDetail.ClassUuid,         //おしらせ引用UUID
		RefUuid:           noticeDetail.RefUuid,           //おしらせ引用UUID
	}

	//確認しているか取得
	status, err := model.GetNoticeReadStatus(noticeUuid, userUuid)
	if err != nil {
		return NoticeDetail{}, err
	}

	fmt.Println(status)
	//確認していた場合、ReadStatusに1を保存する
	formattedNotice.ReadStatus = 0
	if status {
		formattedNotice.ReadStatus = 1
	}

	//userUuidをuserNameに整形
	teacherUuid := noticeDetail.UserUuid
	teacher, nil := model.GetUser(teacherUuid)
	if err != nil {
		return NoticeDetail{}, err
	}
	//整形後formatに追加
	formattedNotice.UserName = teacher.UserName // おしらせ発行ユーザ

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
	NoticeUuid  string    // おしらせUUID
	NoticeTitle string    //お知らせのタイトル
	NoticeDate  time.Time //お知らせの作成日時
	UserName    string    // おしらせ発行ユーザ
	ClassUuid   string    // クラスUUID
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
			NoticeUuid:  notice.NoticeUuid,  //おしらせUuid
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
		notices.ClassUuid = classUuid       // おしらせUuid
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
	isParent, err := model.IsParent(bRead.UserUuid)
	if err != nil { // エラーハンドル
		return err
	}
	if !isParent { // 非管理者ユーザーの場合
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


// 特定のお知らせ既読済み一覧
type NoticeStatus struct {
	UserName   string // ガキの名前
	GenderCode string //性別コード
}

func (s *NoticeService) GetNoticeStatus(noticeUuid string, userUuid string) (NoticeStatus, error) {

	//TODO:削除
	fmt.Println("さーびすです")
	fmt.Println(noticeUuid)

	//おしらせがどのクラスのものなのかを取ってくる
	temp, err := model.GetNoticeDetail(noticeUuid)
	if err != nil {
		return NoticeStatus{}, err
	}
	//わかりやすよう、tempのclassUuidだけ取ってきとく
	classUuid := temp.ClassUuid
	fmt.Println("classUuid:::" + classUuid)

	//お知らせの既読済情報一覧を取ってくる
	noticeReadStatus, err := model.GetNoticeStatusList(noticeUuid)
	if err != nil {
		return NoticeStatus{}, err
	}

	//構造体からスライスに変換する(userUuidを持ってる配列を作る)
	var userUuids []string
	for _, statusList := range noticeReadStatus {
		userUuid := statusList.UserUuid
		userUuids = append(userUuids, userUuid)
	}

	// userUuidsはお知らせを既読しているユーザの一覧を保持
	fmt.Println(userUuids)

	// userUuidからouhciUuidの一覧を作る
	var ouchiUuids []string
	for _, usersList := range userUuids {
		//usesrを一つ取って
		userUuid := usersList
		//getUserでユーザ情報取ってくる
		userDetail, err := model.GetUser(userUuid) // user情報のすべてが返るのでdetailにしてる
		if err != nil {
			return NoticeStatus{}, err
		}

		// ouchiUuidがnullやった場合の処理
		if userDetail.OuchiUuid == nil {
			logging.ErrorLog("OuchiUuid is nil for user UUID %s", err)
			return NoticeStatus{}, err
		}

		//取ってきた情報のouchiUuidを追加していく
		ouchiUuids = append(ouchiUuids, *userDetail.OuchiUuid)
	}

	// ouchiUuid 一覧を保持
	fmt.Println(ouchiUuids)

	noticeStatus := NoticeStatus{
		UserName:   userUuid,
		GenderCode: "0",
	}

	return noticeStatus, err
}

