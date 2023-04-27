package dialect

import "reflect"

/*
隔离不同数据库之间的差异，便于扩展
*/

var dialectMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf 将go语言的类型转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL 返回某个表是否存在的sql语句
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册 Dialect 实例
func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

// GetDialect 获取 Dialect 实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
