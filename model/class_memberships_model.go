package model

// ユーザのクラス所属中間テーブル
type ClassMembership struct {
	ClassUuid string `xorm:"varchar(36) pk" json:"classUUID"` // クラスID
	UserUuid  string `xorm:"varchar(36) pk" json:"userUUID"`  // ユーザーID
}

// テーブル名
func (ClassMembership) TableName() string {
	return "class_memberships"
}

// FK制約の追加
func InitClassMembershipFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE class_memberships ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	// ClassUuid
	_, err = db.Exec("ALTER TABLE class_memberships ADD FOREIGN KEY (class_uuid) REFERENCES classes(class_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}

	return nil
}

// テストデータ
func CreateClassMembershipsTestData() {
	classMemberships1 := &ClassMembership{
		UserUuid:  "3cac1684-c1e0-47ae-92fd-6d7959759224",
		ClassUuid: "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(classMemberships1)
	classMemberships2 := &ClassMembership{
		UserUuid:  "9efeb117-1a34-4012-b57c-7f1a4033adb9",
		ClassUuid: "09eba495-fe09-4f54-a856-9bea9536b661",
	}
	db.Insert(classMemberships2)
	classMemberships3 := &ClassMembership{
		UserUuid:  "3cac1684-c1e0-47ae-92fd-6d7959759224",
		ClassUuid: "817f600e-3109-47d7-ad8c-18b9d7dbdf8b",
	}
	db.Insert(classMemberships3)
}
// user_uuidで絞り込み、所属クラスの構造体のスライスとerrorを返す
func FindClassMemberships(userUuid string) ([]ClassMembership, error) {
	//ClassMemberships構造体のスライスを返すので定義
	var classMemberships []ClassMembership

	//uuidをWhere句で条件指定
	err := db.Where("user_uuid = ?", userUuid).Find(&classMemberships)
	if err != nil { //エラーハンドル
		return nil, err
	}

	return classMemberships, nil
}

// user_uuidで絞り込み、所属クラスの構造体のスライスとerrorを返す
func GetClassList(userUuids []string) ([]ClassMembership, error) {
	//Class構造体のスライスを返すので定義
	var classMemberships []ClassMembership

	// uuidで絞り込み
	err := db.In("user_uuid", userUuids).Find(&classMemberships)
	if err != nil { //エラーハンドル
		return nil, err
	}

	return classMemberships, nil
}

// ユーザーをクラスに所属させるよ
// 新しい構造体をレコードとして受け取り、usersテーブルにinsertし、可否とerrorを返す
func JoinClass(record ClassMembership) (bool, error) {
	affected, err := db.Insert(record)
	if err != nil || affected == 0 { //エラーハンドル
		return false, err // 受け取ったエラーを返す
	}
	return true, nil
}

// クラスIDから参加しているユーザーを全取得
func FindClassMembers(classUuid string) ([]ClassMembership, error) {
	var classMemberships []ClassMembership
	//uuidをWhere句で条件指定
	err := db.Where("class_uuid = ?", classUuid).Find(&classMemberships)
	if err != nil { //エラーハンドル
		return nil, err
	}
	return classMemberships, nil
 }

//クラスに所属しているおこさんだけを全件取得(先生は除外するためuserUuidでnot in)
func FindUserByClassMemberships(classUuid string, userUuid string) ([]ClassMembership, error) {
	//ClassMembership型で返す(あってるのかは知らん「)
	var user []ClassMembership
	//classuuidで絞り込み
	err := db.Where("class_uuid = ?", classUuid).
				Where("user_uuid NOT IN (?)", userUuid).Find(&user)
	if err != nil { //エラーハンドル
		return nil, err
	}

	return user, nil
}
func CheckClassMemberships(userUuids []string, classUuids []string) ([]ClassMembership, error) {
	var  result []ClassMembership
	err := db.In("user_uuid", userUuids).In("class_uuid", classUuids).Find(&result)
	if err != nil {
		return nil, err
	}
	return  result, nil
}