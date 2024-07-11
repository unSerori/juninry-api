package service

import (
	"crypto/rand"
	"fmt"
	"juninry-api/common"
	"juninry-api/logging"
	"juninry-api/model"
	"math/big"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type ClassService struct{}

// 招待コード生成部分
// クラス内でしか呼び出されない
func (s *ClassService) generateInviteCode(bClass model.Class) (model.Class, error) {
	// 有効な招待コードが無ければ新しい招待コードを作る
	// 有効期限を1週間後に設定
	validUntil := time.Now().AddDate(0, 0, 7)
	bClass.ValidUntil = validUntil // バインド

	// 10回エラー吐いたら終わりでええやろ。。。
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		// 招待コードを作る
		inviteCode, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil { // 乱数生成でエラーが出たら泣く
			continue
		}
		// 4桁文字列にキャストしてバインド
		bClass.InviteCode = fmt.Sprintf("%04d", inviteCode.Int64())

		// クラステーブルに追加
		_, err = model.UpdateInviteCode(bClass)
		if err != nil {
			// 本処理時のエラーごとに処理(:DBエラー)
			switch err.(*mysql.MySQLError).Number {
			case 1062: // 一意性制約違反の場合、招待コードから再生成
				continue
			default: // 一意性制約違反じゃなかったらびっくり
				return model.Class{}, err
			}
		}

		// 招待コードを作成できたクラスを返す
		return bClass, nil
	}

	// 試行回数10回以上で失敗したらエラーを返す
	// これ10回連続衝突する可能性そこそこあるよね〜
	// TODO: 改善の余地あり
	logging.ErrorLog("Maximum number of attempts reached", nil)
	return model.Class{}, common.NewErr(common.ErrTypeMaxAttemptsReached)
}

func (s *ClassService) PermissionCheckedClassCreation(userUuid string, bClass model.Class) (model.Class, error) {

	// クラス作成権限を持っているか確認
	isTeacher, err := model.IsTeacher(userUuid)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return model.Class{}, common.NewErr(common.ErrTypePermissionDenied)
	}

	// クラス作成
	// クラスUUIDの生成
	classUuid, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {                    // 空の構造体とエラー
		return model.Class{}, err
	}

	bClass.ClassUuid = classUuid.String() // バインド

	// クラス作成
	_, err = model.CreateClass(bClass)
	if err != nil { // エラーハンドル
		return model.Class{}, err // uuidの作成がおかしくなければ問題ないけど、登録結果が0件で正常終了することなんかあるか？
	}

	// クラスに参加する
	classMembership := model.ClassMembership{
		UserUuid:  userUuid,
		ClassUuid: bClass.ClassUuid,
	}

	// クラスに参加
	success, err := model.JoinClass(classMembership)
	if err != nil || !success { // エラーハンドル
		return model.Class{}, err
	}

	// 招待コード入ったクラスもらえます！
	class, err := s.generateInviteCode(bClass)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}

	//エラーが出なかった場合、コミットして作成したクラスを返す
	return class, nil
}

// 児童の構造体
type JuniorData struct {
	// ここに出席番号を追加
	UserName string `json:"userName"`
	GenderId int    `json:"genderId"`
}
// クラスごとに児童のデータをまとめた構造体
type TransFormData struct {
	ClassName string       `json:"className"`
	JuniorData []JuniorData `json:"juniorData"`
}

// 自身が参加しているクラスに参加しているユーザ情報を全権取得
func (s *ClassService) GetClassMates(useruuid string) ([]TransFormData, error) {

	// useridでユーザ情報を取得
	myClass, err := model.FindClassMemberships(useruuid)
	if err != nil {
		return nil, err
	}
	// for文でクラスのIDを配列に格納
	var classUUIDs []string
	for _,i := range myClass {
		classUUIDs = append(classUUIDs, i.ClassUuid)
	}
	// UUIDにあわせてクラス名を取得
	// クラス名をキー、バリューをデータのマップにする
	transformedDataMap := make(map[string][]JuniorData)
	for _,uuid := range classUUIDs {
		// uuid
		class,err := model.GetClass(uuid)
		if err != nil {
			return nil, err
		}
		// 参加しているユーザーを全取得
		memberships,err := model.FindClassMembers(uuid)
		if err != nil {
			return nil, err
		}
		// ユーザーIDの配列に格納
		var membershipsUUIDs []string
		for _,i := range memberships {
			membershipsUUIDs = append(membershipsUUIDs, i.UserUuid)
		}
		// 配列からユーザー情報を取得
		classmates,err := model.GetUsers(membershipsUUIDs)
		if err != nil {
			return nil, err
		}
		for _, classmate := range classmates {
			juniorData := JuniorData{
				UserName: classmate.UserName,
				GenderId: classmate.GenderId,
			}
			transformedDataMap[class.ClassName] = append(transformedDataMap[class.ClassName], juniorData)
		}
	}
	// つくったマップをさらに成形
	var transformedDataList []TransFormData
	for className, juniorData := range transformedDataMap {
		transformedData := TransFormData{
			ClassName: className,
			JuniorData: juniorData,
		}
		transformedDataList = append(transformedDataList, transformedData)
	}
	return transformedDataList, err
}

func (s *ClassService) PermissionCheckedRefreshInviteCode(userUuid string, classUuid string) (model.Class, error) {

	// クラス作成権限を持っているか確認
	isTeacher, err := model.IsTeacher(userUuid)
	if err != nil { // エラーハンドル
		return model.Class{}, err // トークンあるのにユーザーがいないことはあり得ないのでないと思うが、、、？
	}
	if !isTeacher { // 非管理者ユーザーの場合
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return model.Class{}, common.NewErr(common.ErrTypePermissionDenied)
	}
	// クラスUUIDが存在するかどうか
	targetClass, err := model.GetClass(classUuid)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}
	if targetClass.ClassUuid == "" { // そんなクラス存在しない場合
		return model.Class{}, common.NewErr(common.ErrTypeNoResourceExist)
	}

	// 招待コード入ったクラスもらえます！
	class, err := s.generateInviteCode(targetClass)
	if err != nil { // エラーハンドル
		return model.Class{}, err
	}
	//エラーが出なかった場合、コミットして作成したクラスを返す
	return class, nil
}

// クラスに参加させる
func (s *ClassService) PermissionCheckedJoinClass(userUuid string, inviteCode string) (string, error) {

	// クラス参加権限を持っているか確認
	isPatron, err := model.IsPatron(userUuid)
	if err != nil { // エラーハンドル
		return "", err
	}
	if isPatron { // 親がクラスに直接入ってくるなってやつです
		logging.ErrorLog("Do not have the necessary permissions", nil)
		return "", common.NewErr(common.ErrTypePermissionDenied)
	}

	// クラスUUIDが存在するかどうか
	targetClass, err := model.GetClassByInviteCode(inviteCode)
	if err != nil { // エラーハンドル

		return "", err
	}
	if targetClass.ClassUuid == "" { // そんなクラス存在しない場合
		return "", common.NewErr(common.ErrTypeNoResourceExist)
	}

	// クラスに参加させる
	classMembership := model.ClassMembership{
		UserUuid:  userUuid,
		ClassUuid: targetClass.ClassUuid,
	}

	// クラスに所属しようね
	success, err := model.JoinClass(classMembership)
	if err != nil || !success { // エラーハンドル
		return "", err
	}

	// 所属したクラス名を返す
	class, err := model.GetClass(classMembership.ClassUuid)
	if err != nil { // エラーハンドル
		return "", err
	}

	// エラーが出なかった場合、クラス名を返す
	return class.ClassName, nil
}
