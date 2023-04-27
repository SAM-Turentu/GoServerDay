package test

import (
	"geeorm/log"
	"geeorm/session"
	"testing"
)

// Account 创建表 字段首字母不能小写
type Account struct {
	Id       int `geeorm:"primary key"`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	log.Info("before insert", account)
	account.Id += 100
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", account)
	account.Password = "*****"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123"}, &Account{2, "asdf"})

	u := &Account{} // 返回第一条记录的account password=‘*****’

	err := s.First(u)
	if err != nil || u.Id != 101 || u.Password != "*****" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}
