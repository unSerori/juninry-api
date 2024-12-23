package model

// 教科テーブル
type Subject struct {
	SubjectId   int    `xorm:"pk autoincr" json:"subjectID"`                   // 教科ID
	SubjectName string `xorm:"varchar(15) not null unique" json:"subjectName"` // 強化名
}

// テーブル名
func (Subject) TableName() string {
	return "subjects"
}

// テストデータ
func CreateSubjectTestData() {
	subject1 := &Subject{
		SubjectName: "国語",
	}
	db.Insert(subject1)
	subject2 := &Subject{
		SubjectName: "算数",
	}
	db.Insert(subject2)
	subject3 := &Subject{
		SubjectName: "理科",
	}
	db.Insert(subject3)
	subject4 := &Subject{
		SubjectName: "社会",
	}
	db.Insert(subject4)
	subject5 := &Subject{
		SubjectName: "英語",
	}
	db.Insert(subject5)
}
