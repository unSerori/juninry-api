// エンティティ

package domain

// エンティティとその値オブジェクト
type TeachingMaterial struct {
	TeachingMaterialUuid      string // 教材ID
	TeachingMaterialName      string `json:"teachingMaterialName" form:"teachingMaterialName"` // 教材名
	SubjectId                 int    `json:"subjectId" form:"subjectId"`                       // 教科ID
	TeachingMaterialImageUuid string // 教材画像
	ClassUuid                 string `json:"classUUID" form:"classUUID"` // クラスID
}

// ファクトリー関数
// func NewTeachingMaterial() *TeachingMaterial {
// 	return &TeachingMaterial{}
// }
