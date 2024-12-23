package service

import (
	"errors"
	"fmt"
	"juninry-api/common/custom"
	"juninry-api/common/logging"
	"juninry-api/model"
	"time"

	"github.com/go-sql-driver/mysql"
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
		return custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// notice_uuidを生成
	noticeId, err := uuid.NewRandom() //新しいuuidの作成
	if err != nil {
		return err
	}
	bNotice.NoticeUuid = noticeId.String() //設定

	// 投稿時刻を設定
	bNotice.NoticeDate = time.Now()

	// 構造体をレコード登録処理に投げる
	_, err = model.CreateNotice(bNotice) // 第一返り血は登録成功したレコード数
	if err != nil {                      // エラーハンドル
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
		return err
	}

	return nil

}

// おしらせテーブル(1件取得用)
type NoticeDetail struct { // typeで型の定義, structは構造体
	NoticeUuid        string    `json:"noticeUUID"`        // お知らせUUID
	NoticeTitle       string    `json:"noticeTitle"`       //お知らせのタイトル
	NoticeExplanatory string    `json:"noticeExplanatory"` //お知らせの内容
	NoticeDate        time.Time `json:"noticeDate"`        //お知らせの作成日時
	UserName          string    `json:"userName"`          // おしらせ発行ユーザ
	ClassName         string    `json:"className"`         // どのクラスのお知らせか
	ClassUuid         string    `json:"classUUID"`         // クラスUUID
	QuotedNoticeUuid  *string   `json:"quotedNoticeUUID"`  // 引用お知らせUUID
	ReadStatus        *int      `json:"readStatus"`        // 既読ステータス
}

// お知らせ詳細取得
func (s *NoticeService) GetNoticeDetail(noticeUuid string, userUuid string) (NoticeDetail, error) {
	// お知らせを確認する権限があるか確認
	user, err := model.GetUser(userUuid)
	if err != nil {
		return NoticeDetail{}, err
	}

	var ouchiUuid string
	if user.OuchiUuid != nil { // あったらいいな、お家（なくても既読見えないだけだから別に関係ない）
		ouchiUuid = *user.OuchiUuid
	}

	// 閲覧許可のあるクラス一覧
	allowedClassUuids := []string{userUuid}

	isPatron, err := model.IsPatron(userUuid)
	if err != nil {
		return NoticeDetail{}, err
	}

	// 親の場合は子供をたどり、そうでない場合は自身の所属クラスを確認
	if isPatron {
		// 親の場合お家IDがなかったら話にならないので破壊
		if ouchiUuid == "" {
			return NoticeDetail{}, custom.NewErr(custom.ErrTypeNoResourceExist)
		}

		// 同じお家IDの子供のユーザーIDを取得
		userUuids, err := model.GetChildrenUuids(*user.OuchiUuid)
		if err != nil { // エラーハンドル
			return NoticeDetail{}, err
		}

		if len(userUuids) == 0 {
			//  エラー:おうちに子供はいないのになにしてんのエラー
			return NoticeDetail{}, custom.NewErr(custom.ErrTypeNoResourceExist)
		}

		// 子供のクラスUUID一覧取得
		classes, err := model.GetClassList(userUuids)
		if err != nil { // エラーハンドル
			return NoticeDetail{}, err
		}

		// 親が閲覧許可のあるクラスたち
		for _, class := range classes {
			allowedClassUuids = append(allowedClassUuids, class.ClassUuid)
		}
	} else {
		// 自身の閲覧可能クラス
		classes, err := model.FindClassMemberships(userUuid)
		if err != nil { // エラーハンドル
			return NoticeDetail{}, err
		}
		for _, class := range classes {
			allowedClassUuids = append(allowedClassUuids, class.ClassUuid)
		}
	}

	//お知らせ詳細情報取得
	noticeDetail, err := model.GetNoticeDetail(noticeUuid)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return NoticeDetail{}, err //nilで返せない!不思議!!  // A. 返り血の方がNoticeDetailになっていてNoticeDetail型で返さなければいけないから。*NoticeDetailのようにポインタで返せばポインタの指定先が空の状態≒nilを返すことができるよ。
	}

	if noticeDetail == nil { // 取得できなかった
		fmt.Println("noticeDetail is nil")
		return NoticeDetail{}, custom.NewErr(custom.ErrTypeNoResourceExist)
	}

	// 取得したお知らせにアクセスする許可があるか確認
	found := false
	for _, uuid := range allowedClassUuids {
		if uuid == noticeDetail.ClassUuid {
			found = true
			break
		}
	}

	if !found {
		return NoticeDetail{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	//取ってきたnoticeDetailを整形して、controllerに返すformatに追加する
	formattedNotice := NoticeDetail{
		NoticeUuid:        noticeDetail.NoticeUuid,        // お知らせUUID
		NoticeTitle:       noticeDetail.NoticeTitle,       // お知らせタイトル
		NoticeExplanatory: noticeDetail.NoticeExplanatory, // お知らせの内容
		NoticeDate:        noticeDetail.NoticeDate,        // お知らせ作成日時
		QuotedNoticeUuid:  noticeDetail.QuotedNoticeUuid,  // 引用おしらせUUID
		ClassUuid:         noticeDetail.ClassUuid,         // クラスUUID
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

	// お家に所属している場合、お知らせの既読状況確認
	if ouchiUuid != "" {
		//確認しているか取得
		status, err := model.IsRead(noticeUuid, ouchiUuid)
		if err != nil {
			return NoticeDetail{}, err
		}

		//確認していた場合、ReadStatusに1を保存する
		if status {
			read := 1
			formattedNotice.ReadStatus = &read
		} else {
			unRead := 0
			formattedNotice.ReadStatus = &unRead
		}
	}
	return formattedNotice, err
}

// 子ども情報テーブル
type PupilInfo struct { // typeで型の定義, structは構造体
	PupilUuid string `json:"pupilUUID"` // こどものUUID
	PupilName string `json:"pupilName"` // こどもの名前
}

// おしらせテーブル(全件取得用)
type NoticeHeader struct { // typeで型の定義, structは構造体
	NoticeUuid  string      `json:"noticeUUID"`  // おしらせUUID
	NoticeTitle string      `json:"noticeTitle"` //お知らせのタイトル
	NoticeDate  time.Time   `json:"noticeDate"`  //お知らせの作成日時
	UserName    string      `json:"userName"`    // おしらせ発行ユーザ
	ClassUuid   string      `json:"classUUID"`   // クラスUUID
	ClassName   string      `json:"className"`   // どのクラスのお知らせか
	ReadStatus  *int        `json:"readStatus"`  //お知らせを確認しているか　お家に所属していない場合、既読情報は存在しないのでポインタを使いnilを許容する
	PupilInfos  []PupilInfo `json:"pupilInfo"`   // お知らせのクラスに所属している子供一覧
}

// ユーザの所属するクラスのお知らせ全件取得
func (s *NoticeService) FindAllNotices(userUuid string, classUuids []string, pupilUuids []string, sortReadStatus *int) ([]NoticeHeader, error) {

	// 結果格納用変数
	var userUuids []string

	// 既読テーブル問い合わせに必要なお家UUIDを定義
	var ouchiUuid string

	// とりあえずユーザーテーブルからもろもろを引っ張ってくる
	user, err := model.GetUser(userUuid)
	if err != nil { // エラーハンドル
		return nil, err
	}
	if user.OuchiUuid != nil { // あったらいいな、お家（なくても既読見えないだけだから別に関係ない）
		ouchiUuid = *user.OuchiUuid
	}

	isPatron, err := model.IsPatron(userUuid)
	if err != nil { // エラーハンドル
		return nil, err
	}

	if isPatron { // ユーザーが親の場合は子供のIDを取得する必要
		// 親の場合お家IDがなかったら話にならないので破壊
		if ouchiUuid == "" {
			return nil, custom.NewErr(custom.ErrTypeNoResourceExist)
		}

		// 絞り込み条件がなければ、同じおうちの子供全員を取得してくる
		if len(pupilUuids) == 0 {
			// 同じお家IDの子供のユーザーIDを取得
			userUuids, err = model.GetChildrenUuids(*user.OuchiUuid)
			if err != nil { // エラーハンドル
				return nil, err
			}
		} else { // pupilUuidsがあれば同じおうちに所属しているのか調べる
			// Uuidすべてを検索する
			for _, pupilUuid := range pupilUuids {
				pupil, err := model.GetUser(pupilUuid)
				if err != nil {
					return nil, err
				}

				fmt.Println("親", *user.OuchiUuid, ":子供", *pupil.OuchiUuid)
				// 親のouchiUuidと比べて違うかったらえらーよね
				if *user.OuchiUuid != *pupil.OuchiUuid {
					// そんな子供いないよエラー
					return nil, custom.NewErr(custom.ErrTypeNoResourceExist)
				}
			}
			// 絞込みの子供IDを保存(kisyoi,dogeza)
			userUuids = pupilUuids
		}

		if len(userUuids) == 0 {
			//  エラー:おうちに子供はいないのになにしてんのエラー
			return nil, custom.NewErr(custom.ErrTypeNoResourceExist)
		}

	} else { // 親でない場合はそのまま自分のIDでと合わせることができる
		userUuids = append(userUuids, userUuid)
	}

	// お知らせの取得条件に使うclassUUIDのスライス
	var displayedClassUuids []string

	// classUuidsが空の場合の処理(絞り込みなしの全件取得)
	if len(classUuids) == 0 {

		// userUuidを条件にしてclassUuidを取ってくる
		// 1 - userUuidからclass_membershipの構造体を取ってくる
		classMemberships, err := model.GetClassList(userUuids)
		if err != nil {
			return nil, err
		}

		// 2 - 構造体からclassUuidのスライス(配列)を作る
		for _, classMembership := range classMemberships {
			classUuid := classMembership.ClassUuid
			displayedClassUuids = append(displayedClassUuids, classUuid)
		}

	} else { // classUuidで絞り込まれた取得(絞り込み条件がなかったらエラーだよネ)
		// ユーザーがクラスに所属しているかを確認する
		classMemberships, err := model.CheckClassMemberships(userUuids, classUuids)

		if err != nil || classMemberships == nil {
			logging.ErrorLog("Do not have the necessary permissions", nil)
			return []NoticeHeader{}, custom.NewErr(custom.ErrTypePermissionDenied)
		}

		// 2 - 構造体からclassUuidのスライス(配列)を作る
		for _, classMembership := range classMemberships {
			classUuid := classMembership.ClassUuid
			displayedClassUuids = append(displayedClassUuids, classUuid)
		}
	}

	// classUuidを条件にしてnoticeの構造体を取ってくる
	notices, err := model.FindNotices(displayedClassUuids)
	if err != nil {
		return nil, err
	}

	//format後のnotices格納用変数(複数返ってくるのでスライス)
	var noticeHeaders []NoticeHeader

	//noticesの一つをnoticeに格納(for文なのでデータ分繰り返す)
	for _, notice := range notices {

		// そのまま挿入できるデータを突っ込む
		noticeHeader := NoticeHeader{
			NoticeUuid:  notice.NoticeUuid,  //おしらせUuid
			NoticeTitle: notice.NoticeTitle, //お知らせのタイトル
			NoticeDate:  notice.NoticeDate,  //お知らせの作成日時
			ClassUuid:   notice.ClassUuid,   //お知らせのクラスUuid
		}

		// お知らせ作成者の名前をUUIDから取得する
		creatorUser, err := model.GetUser(notice.UserUuid) //ユーザ取得
		if err != nil {
			return []NoticeHeader{}, err
		}
		//整形後formatに追加
		noticeHeader.UserName = creatorUser.UserName // おしらせ発行ユーザ

		// お知らせのクラス名をUUIDから取得する
		classUuid := notice.ClassUuid
		class, nil := model.GetClass(classUuid) //クラス取得
		if err != nil {
			return []NoticeHeader{}, err
		}
		//整形後formatに追加
		noticeHeader.ClassName = class.ClassName // おしらせ発行ユーザ

		//classUuidから子どもを特定する
		pupils, err := model.GetUserByClassUuid(classUuid, userUuids)
		if err != nil {
			return []NoticeHeader{}, err
		}

		fmt.Println(pupils)
		//返す値を格納する変数
		var pupilInfos []PupilInfo
		//infoに情報を入れていく
		for _, pupil := range pupils {
			user, err := model.GetUser(pupil.UserUuid)
			if err != nil {
				return []NoticeHeader{}, err
			}
			// 必要な値だけ取り出して代入
			pupilInfo := PupilInfo{
				PupilUuid: user.UserUuid,
				PupilName: user.UserName,
			}
			// スライスに追加
			pupilInfos = append(pupilInfos, pupilInfo)
		}

		//宣言した構造体に情報をいれる
		noticeHeader.PupilInfos = pupilInfos

		// リクエストしたユーザーがお家を持っている場合既読ステータスを取得
		if ouchiUuid != "" {
			// 確認しているか取得
			status, err := model.IsRead(notice.NoticeUuid, ouchiUuid)
			if err != nil {
				return []NoticeHeader{}, err
			}

			//確認していた場合、ReadStatusに1を保存する
			if status {
				Read := 1
				noticeHeader.ReadStatus = &Read
			} else {
				Unread := 0
				noticeHeader.ReadStatus = &Unread
			}
		}

		// 直接nilとの比較ができないので条件がない時の値を定義nil
		var noFilterStatus *int

		if sortReadStatus != noFilterStatus {

			// XXX:絞り込み条件が予測しない値の場合弾く
			if *sortReadStatus != 0 && *sortReadStatus != 1 {
				logging.ErrorLog("Do not have the necessary permissions", nil)
				return []NoticeHeader{}, custom.NewErr(custom.ErrTypeUnforeseenCircumstances)
			}

			// ソート条件があるなら、条件に合うものだけ追加していく
			if *sortReadStatus == *noticeHeader.ReadStatus { //　条件在り(0か1)
				noticeHeaders = append(noticeHeaders, noticeHeader)
			}
		} else { //　条件なし
			noticeHeaders = append(noticeHeaders, noticeHeader)
		}
	}

	return noticeHeaders, nil
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
		return custom.NewErr(custom.ErrTypePermissionDenied)
	}

	user, err := model.GetUser(userUuid)
	if err != nil {
		return err
	}
	fmt.Println("今からレコード登録")

	readStatus := model.NoticeReadStatus{
		NoticeUuid: noticeUuid,
		OuchiUuid:  *user.OuchiUuid,
	}

	// 構造体をレコード登録処理に投げる
	err = model.ReadNotice(readStatus) // 第一返り血は登録成功したレコード数
	if err != nil {
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
		return err
	}

	return nil
}

// 特定のお知らせ既読済み一覧 TODO:出席番号どうする？
type NoticeStatus struct {
	StudentNumber *int   `json:"studentNumber"` // 出席番号
	UserName      string `json:"userName"`      // ガキの名前
	GenderId      *int   `json:"genderId"`      // 性別コード(定義がないためnullにしてる)
	ReadStatus    *int   `json:"readStatus"`    // お知らせを確認しているか
}

// 特定のお知らせ既読済み一覧取得
func (s *NoticeService) GetNoticeStatus(noticeUuid string, userUuid string) ([]NoticeStatus, error) {

	// 取得権限を持っているか確認
	isTeacher, err := model.IsTeacher(userUuid)
	if err != nil { // エラーハンドル
		return nil, err
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return nil, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	//確認用です
	fmt.Println("noticeUuid:"+noticeUuid, "userUuid:"+userUuid)

	//おしらせがどのクラスのものなのかを取ってくる
	notice, err := model.GetNoticeDetail(noticeUuid)
	if err != nil {
		return []NoticeStatus{}, err
	}
	//わかりやすよう、noticeのclassUuidだけ取ってきとく
	classUuid := notice.ClassUuid

	// TODO:　教員が複数人いた時おかしくなりますよね
	classMemberships, err := model.FindUserByClassMemberships(classUuid, userUuid)
	if err != nil {
		return []NoticeStatus{}, err
	}

	// userUUIDをキーとしたマップを作成
	studentList := make(map[string]NoticeStatus)
	for _, membership := range classMemberships {
		studentList[membership.UserUuid] = NoticeStatus{
			StudentNumber: membership.StudentNumber,
		}
	}

	// 不足している項目をuserテーブルから取得
	for userUuid, student := range studentList {
		user, err := model.GetUser(userUuid)
		if err != nil {
			return []NoticeStatus{}, err
		}
		// student.GenderId = &user.GenderId			// 性別コード
		student.UserName = user.UserName // 名前を挿入
		// 既読状況の取得
		if user.OuchiUuid != nil {
			result, err := model.IsRead(noticeUuid, *user.OuchiUuid)
			if err != nil {
				return []NoticeStatus{}, err
			}
			if result { // 既読
				read := 1
				student.ReadStatus = &read
			} else {
				unRead := 0
				student.ReadStatus = &unRead
			}
		}

		studentList[userUuid] = student // 更新したものを入れ直す
	}

	// Mapをスライスに変換
	// HACK: 生徒0人だとnullが帰ってきてうざいのでからの長さ0のスライスを作成
	noticeStatusList := []NoticeStatus{}
	for _, student := range studentList {
		noticeStatusList = append(noticeStatusList, student)
	}
	return noticeStatusList, nil
}
