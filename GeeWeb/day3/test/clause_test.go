package test

import (
	clause2 "geeorm/clause"
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	var clause clause2.Clause
	clause.Set(clause2.LIMIT, 3)
	clause.Set(clause2.SELECT, "User", []string{"*"})
	clause.Set(clause2.WHERE, "Name = ?", "Tom")
	clause.Set(clause2.ORDERBY, "Age asc")
	sql, vars := clause.Build(clause2.SELECT, clause2.WHERE, clause2.ORDERBY, clause2.LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age asc LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}
