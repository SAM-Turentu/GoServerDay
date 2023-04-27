package session

import (
	"database/sql"
	"geeorm/dialect"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

//engine, _ := geeorm.NewEngine("sqlite3", "D:\\Common\\Sqlite\\db\\gee.db")

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "D:\\Common\\Sqlite\\db\\gee.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	reslut, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "TOm", "SAM").Exec()
	if count, err := reslut.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, bug got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("select count(*) form User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}
