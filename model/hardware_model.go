package model

import (
	"errors"
	"juninry-api/common/custom"

	"github.com/go-sql-driver/mysql"
)

// ハードウェアテーブル

type Hardware struct {
	HardwareUuid   string `xorm:"varchar(36) pk" json:"hardwareUUID"`
	HardwareTypeId int    `xorm:"int not null" json:"hardwareTypeId"`
}

// テーブル名
func (Hardware) TableName() string {
	return "hardwares"
}

// FK制約の追加
func InitHardwareFK() error {
	_, err := db.Exec("ALTER TABLE hardwares ADD FOREIGN KEY (hardware_type_id) REFERENCES hardware_types(hardware_type_id) ON DELETE CASCADE ON UPDATE CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// テスト用
func CreateHardwareTestData() {
	hardware1 := &Hardware{
		HardwareUuid:   "df2b1f4c-b49a-4068-80c5-3120dceb14c8",
		HardwareTypeId: 1,
	}
	db.Insert(hardware1)

	hardware2 := &Hardware{
		HardwareUuid:   "d611d471-5eb2-46a2-abaf-f758205f0d5f",
		HardwareTypeId: 1,
	}
	db.Insert(hardware2)
}

// 新規登録
func CreateHardware(record Hardware) error {
	_, err := db.Insert(record)
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
		// 通常の処理エラー
		return err
	}

	return nil
}
