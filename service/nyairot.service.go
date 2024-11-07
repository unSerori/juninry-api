package service

import (
	"fmt"
	"juninry-api/model"
	"juninry-api/utility/custom"
	"time"
)

type NyariotSarvice struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

// スタンプが増加したかと現在の数を一緒に返すためのテーブル
type StampResult struct {
	StampIncreased bool `json:"stampIncreased"` //　増加したかのture、false
	Quantity       int  `json:"quantity"`       // 現在のスタンプ数
}

// その日初めてのログインの時スタンプを付与
func (s *NyariotSarvice) AddLoginStamp(userUuid string) (StampResult, error) {

	// ユーザーが生徒かな生徒じゃなかったらエラー
	isJunior, err := model.IsJunior(userUuid)
	if err != nil {
		return StampResult{}, err
	}
	if !isJunior {
		return StampResult{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	fmt.Println("生徒チェック終わり")

	// userのスタンプカードを取ってくる
	user, err := model.GetUserStampCard(userUuid)
	if err != nil {
		return StampResult{}, err
	}

	fmt.Println("スタンプカード取得終わり")

	fmt.Println(user)

	// 今の時間を取ってくる
	todayDate := time.Now().Truncate(24 * time.Hour)

	fmt.Println("今日", todayDate)
	fmt.Println("最終ログイン", user.LastLoginTime)

	//　増加したかを保持する変数
	var increased bool

	// 今日が取得日より後かを判定
	if todayDate.YearDay() != user.LastLoginTime.YearDay() {
		// スタンプを増やす
		quantity := user.Quantity + 1

		// スタンプを付与
		_, err = model.AddStamp(userUuid, quantity)
		if err != nil {
			return StampResult{}, err
		}

		// ログイン時間の更新
		_, err = model.UpdateLastLoginTime(userUuid, todayDate)
		if err != nil {
			return StampResult{}, err
		}

		// スタンプが増加したことを保持
		increased = true

	} else {
		// ログイン時間の更新
		_, err = model.UpdateLastLoginTime(userUuid, todayDate)
		if err != nil {
			return StampResult{}, err
		}

		// スタンプの数が変動していないことを保持
		increased = false
	}

	// 更新後のスタンプの数をかえす
	user, err = model.GetUserStampCard(userUuid)
	if err != nil {
		return StampResult{}, err
	}

	Result := StampResult{
		StampIncreased: increased,
		Quantity:       user.Quantity,
	}

	return Result, nil
}

// スタンプの数を取得
func (s *NyariotSarvice) GetStamp(userUuid string) (model.Stamp, error) {
	// ユーザーが生徒かな生徒じゃなかったらエラー
	isJunior, err := model.IsJunior(userUuid)
	if err != nil {
		return model.Stamp{}, err
	}
	if !isJunior {
		return model.Stamp{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// userのスタンプカードを取ってくる
	user, err := model.GetUserStampCard(userUuid)
	if err != nil {
		return model.Stamp{}, err
	}

	fmt.Println(*user)

	return *user, nil
}

// スタンプでガチャ
func (s *NyariotSarvice) GetStampGacha(userUuid string) (model.Item, error) {

	// ユーザーが生徒かな生徒じゃなかったらエラー
	isJunior, err := model.IsJunior(userUuid)
	if err != nil {
		return model.Item{}, err
	}
	if !isJunior {
		return model.Item{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	return model.Item{}, nil
}

// 図鑑用に所持、未所持を保持する変数があるテーブル
type ItemCatalog struct {
	ItemUuid      string `json:"itemUUID"`      // アイテムUUID
	ItemName      string `json:"itemName"`      // アイテム名
	ImagePath     string `json:"imagePath"`     // アイテム画像パス
	ItemNumber    int    `json:"itemNumber"`    // アイテム番号
	Detail        string `json:"detail"`        // アイテム詳細
	Talk          string `json:"talk"`          // アイテム固有の会話
	SatityDegrees int    `json:"satityDegrees"` // 空腹増加値
	Rarity        int    `json:"rarity"`        // アイテムレアリティ 1:N 2:R 3:SR
	HasItem       bool   `json:"hasItem"`       // 所持、未所持
	Quantity      int    `json:"quantity"`      //所持数
}

// 所持アイテム取得
func (s *NyariotSarvice) GetUserItems(userUuid string) ([]ItemCatalog, error) {
	// ユーザーが生徒かな生徒じゃなかったらエラー
	isJunior, err := model.IsJunior(userUuid)
	if err != nil {
		return []ItemCatalog{}, err
	}
	if !isJunior {
		return []ItemCatalog{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// 全アイテムを取得
	items, err := model.GetItems()
	if err != nil {
		return []ItemCatalog{}, err
	}

	// ItemCatalogのスライスを作成
	var itemCatalog []ItemCatalog

	// ユーザが所持しているか確認するためにスライスを作る
	// 同時に返すItemCatalogにデータを入れる
	for _, item := range items {
		// アイテム情報を格納
		catalog := ItemCatalog{
			ItemUuid:      item.ItemUuid,      // アイテムUUID
			ItemName:      item.ItemName,      // アイテム名
			ImagePath:     item.ImagePath,     // アイテム画像パス
			ItemNumber:    item.ItemNumber,    // アイテム番号
			Detail:        item.Detail,        // アイテム詳細
			Talk:          item.Talk,          // アイテム固有の会話
			SatityDegrees: item.SatityDegrees, // 空腹増加値
			Rarity:        item.Rarity,        // アイテムレアリティ 1:N 2:R 3:SR
		}

		// 持ってるか確認
		quantity, has, err := model.GetUserItemBox(userUuid, item.ItemUuid)
		if err != nil {
			return nil, err
		}

		// アイテムを持っていれば、数量とフラグをセット
		if has {
			catalog.HasItem = has       // アイテム持ってる
			catalog.Quantity = quantity // 所持数
		} else {
			catalog.HasItem = has // アイテム持ってない
			catalog.Quantity = 0  // 所持数は0
		}

		// アイテム情報をリストに追加
		itemCatalog = append(itemCatalog, catalog)
	}

	// アイテムリストを返す
	return itemCatalog, nil
}

func (s *NyariotSarvice) GetItemDetail(userUuid string, itemUuid string) (ItemCatalog, error) {

	// ユーザーが生徒かな生徒じゃなかったらエラー
	isJunior, err := model.IsJunior(userUuid)
	if err != nil {
		return ItemCatalog{}, err
	}
	if !isJunior {
		return ItemCatalog{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	// 持ってるか確認
	quantity, has, err := model.GetUserItemBox(userUuid, itemUuid)
	if err != nil {
		return ItemCatalog{}, err
	}

	// 持っていない人は閲覧権限無し
	if !has {
		return ItemCatalog{}, custom.NewErr(custom.ErrTypePermissionDenied)
	}

	item, err := model.GetItem(itemUuid)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ItemCatalog{}, err
	}

	// アイテム情報を格納
	catalog := ItemCatalog{
		ItemUuid:      item.ItemUuid,      // アイテムUUID
		ItemName:      item.ItemName,      // アイテム名
		ImagePath:     item.ImagePath,     // アイテム画像パス
		ItemNumber:    item.ItemNumber,    // アイテム番号
		Detail:        item.Detail,        // アイテム詳細
		Talk:          item.Talk,          // アイテム固有の会話
		SatityDegrees: item.SatityDegrees, // 空腹増加値
		Rarity:        item.Rarity,        // アイテムレアリティ 1:N 2:R 3:SR
		HasItem:       true,               // 所持してるよ
		Quantity:      quantity,           // アイテム個数
	}

	// 詳細を返す
	return catalog, nil

}
