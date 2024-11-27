package model

import (
	"fmt"
	"time"
)

// スタンプテーブル
type Stamp struct {
	UserUuid      string    `xorm:"varchar(36) pk" json:"userUUID"`     // ユーザのUUID
	Quantity      int       `xorm:"int" json:"quantity"`                // スタンプの数
	LastLoginTime time.Time `xorm:"DATE not null" json:"lastLoginTime"` // 最後にログインした時間
}

// テーブル名
func (Stamp) TableName() string {
	return "stamps"
}

// FK制約の追加
func InitStampFK() error {
	// UserUuid
	_, err := db.Exec("ALTER TABLE stamps ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テストデータ
func CreateStampTestData() {
	stamp1 := &Stamp{
		UserUuid:      "3cac1684-c1e0-47ae-92fd-6d7959759224",
		Quantity:      0,
		LastLoginTime: time.Now().Add(-24 * time.Hour),
	}

	db.Insert(stamp1)
}

// スタンプカードの作成
func CreateStampCard(record Stamp) (int64, error) {
	affected, err := db.Insert(record)
	fmt.Printf("Inserted LastGohanTime: %v\n", record.LastLoginTime)
	return affected, err
}

// ユーザのスタンプカードを取得
func GetUserStampCard(userUuid string) (*Stamp, error) {
	//　結果格納用変数
	var stampCard Stamp

	// userUUIDで取得してくる
	found, err := db.Where("user_uuid = ?", userUuid).Get(
		&stampCard,
	)

	// エラーが発生した
	if err != nil {
		return nil, err
	}
	// データが取得できなかったら
	if !found {
		return nil, nil
	}

	return &stampCard, nil
}

// スタンプを追加する
func AddStamp(userUuid string, quantity int) (int64, error) {
	// 更新する構造体のインスタンスを作成
	stamp := Stamp{Quantity: quantity}

	// スタンプ数を更新する
	affected, err := db.Where("user_uuid = ?", userUuid).Cols("quantity").Update(&stamp)
	return affected, err
}

// 　最終ログイン時間の更新
func UpdateLastLoginTime(userUuid string, todayDate time.Time) (int64, error) {
	// 更新する構造体のインスタンスを作成
	stamp := Stamp{LastLoginTime: todayDate}

	// 時間を更新する
	affected, err := db.Where("user_uuid = ?", userUuid).Cols("last_login_time").Update(&stamp)
	return affected, err
}

func ReduceStampQuantity(userUuid string) (int64, error) {

	affected, err := db.Where("user_uuid = ?", userUuid).
		Incr("quantity", -7).
		Update(&Stamp{})
	if err != nil {
		return 0, err
	}
	return affected, err
}
