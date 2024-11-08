package utility

import (
	"fmt"
	"reflect"
)

// 一致するフィールドのみコピーして、構造体を変換する
// 呼び出し側はポインタで渡す ConvertStructCopyMatchingFields(&src, &tgt)
// 埋め込み構造体(:継承された構造体)の場合は中身を再帰的に確認する
func ConvertStructCopyMatchingFields(src interface{}, tgt interface{}) {
	// 値の確認
	fmt.Println("before src struct")
	CheckStruct(src)
	fmt.Println("before tgt struct")
	CheckStruct(tgt)

	// ポインタがさす構造体の実体を取得
	srcVal := reflect.ValueOf(src).Elem()
	tgtVal := reflect.ValueOf(tgt).Elem()

	// フィールドを再帰的にコピー
	CopyMatchingFieldsRecursively(srcVal, tgtVal)

	// 値の確認
	fmt.Println("after tgt struct")
	CheckStruct(tgt)
}

// フィールドコピー処理 もしフィールドが構造体だった場合、それに対して再帰的に処理を行う
func CopyMatchingFieldsRecursively(srcVal reflect.Value, tgtVal reflect.Value) {
	// 元構造体の型
	srcType := srcVal.Type()

	// 構造体のフィールド数だけループ
	for i := 0; i < srcVal.NumField(); i++ {
		// src構造体のi番目のフィールドの型情報と値を取得
		srcField := srcType.Field(i)
		srcFieldValue := srcVal.Field(i)

		// もしフィールドが埋め込み構造体(:継承: フィールドが空で型が埋め込む構造体)なら、再帰的に処理 // 「もしフィールドが構造体なら、再帰的に処理」 としたい場合は srcFieldValue.Kind() == reflect.Struct
		if srcField.Anonymous {
			CopyMatchingFieldsRecursively(srcFieldValue, tgtVal)
		} else {
			// tgt構造体のフィールドたちから、srcのフィールド名と同じものを探す
			if tgtField := tgtVal.FieldByName(srcField.Name); tgtField.IsValid() {
				// 代入可能な型かを確認
				if srcFieldValue.Type().AssignableTo(tgtField.Type()) {
					// 対応するフィールドに設定
					tgtField.Set(srcFieldValue)
				}
			}
		}
	}
}

// 構造体の中身をチェック
func CheckStruct(src interface{}) {
	fmt.Println("CheckStruct: start ============")

	// src構造体の型と値を取得
	sValue := reflect.ValueOf(src)    // 値を取得
	if sValue.Kind() == reflect.Ptr { // 受け取った値がポインタなら、
		sValue = sValue.Elem() // 実際の値を取得
	}
	sType := sValue.Type() // 型を取得

	// src構造体の型と値を取得　// これでは引数が構造体インスタンスor構造体ポインタインスタンスで挙動が変わってしまう
	// sType := reflect.TypeOf(src)  // 型を取得
	// sValue := reflect.ValueOf(src) // 値を取得

	// 構造体のフィールド数だけループ
	for i := 0; i < sType.NumField(); i++ {
		fieldName := sType.Field(i).Name                          // フィールド名を取得
		fieldValue := sValue.Field(i)                             // フィールドの値を取得
		fmt.Printf("%s: %v\n", fieldName, fieldValue.Interface()) // フィールド名と値を出力
	}

	fmt.Println("CheckStruct: end============")
	fmt.Println()
}
