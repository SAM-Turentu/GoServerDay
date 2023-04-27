package test

import (
	"geeorm/session"
	"testing"
)

type UserInfo struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &UserInfo{"Lili", 22}
	user2 = &UserInfo{"zhenzhen", 18}
	user3 = &UserInfo{"aiai", 21}
)

// testRecordInit 初始化表及数据
func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	s := NewSession().Model(&UserInfo{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []UserInfo
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Lili").Update("Age", 19)
	u := &UserInfo{}
	_ = s.OrderBy("Age desc").First(u)
	if affected != 1 || u.Age != 19 {
		t.Fatal("failed to update")
	}
}

func TestSession_Delete(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Lili").Delete()
	count, _ := s.Count()
	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
