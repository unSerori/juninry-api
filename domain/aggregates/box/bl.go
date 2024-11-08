// エンティティのビジネスロジック

package box

/*
箱の状態管理についてのビジネスロジック群
*/

// 数値として取り出す
func (e *Box) GetStatusAsInt() int {
	return e.Status.AsInt()
}

// 更新？
// func (e *Box) ()  {

// }
