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
func (s *NoticeService) ReadNotice(noticeUuid string, userUuid string) error {

	// 既読権限を持っているか確認
	isPatron, err := model.IsPatron(userUuid)
	if err != nil { // エラーハンドル
		return err
	}
	if !isPatron { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return common.NewErr(common.ErrTypePermissionDenied)
	}

	user, err := model.GetUser(userUuid)
	if err != nil {
		return err
	}

	// 構造体をレコード登録処理に投げる
	err = model.ReadNotice(noticeUuid, *user.OuchiUuid) // 第一返り血は登録成功したレコード数
	if err != nil {
		return err
	}

	return nil
}

// 特定のお知らせ既読済み一覧 TODO:出席番号どうする？
type NoticeStatus struct {
	StudentNo  int    // 出席番号
	UserName   string // ガキの名前
	GenderCode string // 性別コード
	ReadStatus int    // お知らせを確認しているか
}

// 特定のお知らせ既読済み一覧取得
func (s *NoticeService) GetNoticeStatus(noticeUuid string, userUuid string) ([]NoticeStatus, error) {

	//TODO:削除
	fmt.Println("さーびすです")
	fmt.Println("noticeUuid：" + noticeUuid)
	fmt.Println("userUuid：" + userUuid)

	//おしらせがどのクラスのものなのかを取ってくる
	notice, err := model.GetNoticeDetail(noticeUuid)
	if err != nil {
		return []NoticeStatus{}, err
	}
	//わかりやすよう、noticeのclassUuidだけ取ってきとく
	classUuid := notice.ClassUuid

	//お知らせの既読済おうち一覧を取ってくる(noticeReadStatus=ouchiuuidみたいなもん)
	noticeReadStatus, err := model.GetNoticeStatusList(noticeUuid)
	if err != nil {
		return []NoticeStatus{}, err
	}
	// 結果をfmtで出力する
	for _, status := range noticeReadStatus {
		fmt.Printf("ouchiUuid: %s", status.OuchiUuid)
	}

	// noticeReadStatusから既読済みガキ一覧を作る
	var readList []model.User
	for _, ouchi := range noticeReadStatus {
		// ouchi.OuchiUuid としてフィールド名を大文字で始める
		gaki, err := model.GetJunior(ouchi.OuchiUuid)
		if err != nil {
			return []NoticeStatus{}, err
		}
		// リストに追加していく
		readList = append(readList, gaki)
	}

	//classUuidからクラス全員を取ってくる(先生は除外するためuserUuidでnotin)
	classMemberships, err := model.FindUserByClassMemberships(classUuid, userUuid)
	if err != nil {
		return []NoticeStatus{}, err
	}

	// もはや、レシピみたいに書いた方がわかりやすいのでわ(脳死)
	// classMembershipsからガキ一覧を作る(ついでに返す奴にデータを突っ込む)
	var juniorList []model.User
	for _, junior := range classMemberships {
		gaki, err := model.GetUser(junior.UserUuid)
		if err != nil {
			return []NoticeStatus{}, err
		}

		// リストに追加していく
		juniorList = append(juniorList, gaki)
	}

	//既読済みガキ一覧でマップを作成
	readMap := make(map[string]bool)
	for _, junior := range readList {
		readMap[junior.UserUuid] = true
	}

	var temp []NoticeStatus

	//juniorlistをループしてマップ検索と整形
	for _, junior := range juniorList {

		// 整形用
		noticeStatus := NoticeStatus{}

		if readMap[junior.UserUuid] {
			noticeStatus.ReadStatus = 1 //　既読済みフラグ
		} else {
			noticeStatus.ReadStatus = 0 //　未読
		}

		noticeStatus.UserName = junior.UserName

		temp = append(temp, noticeStatus)
	}

	//fomat後のnotices格納用変数(複数返ってくるのでスライス)
	var noticeStatus []NoticeStatus

	// tempを逆順にしてformattedAllNoticesに追加する
	for i := len(temp) - 1; i >= 0; i-- {
		noticeStatus = append(noticeStatus, temp[i])
	}

	fmt.Println(noticeStatus)

	return noticeStatus, err
}
