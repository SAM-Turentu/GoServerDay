package main

import (
	"fmt"
	"geeorm"
	_ "github.com/mattn/go-sqlite3"
)

//region main demo
//func main() {
//	db, _ := sql.Open("sqlite3", "D:\\Common\\Sqlite\\db\\gee.db")
//	defer func() { _ = db.Close() }()
//	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
//	//_, _ = db.Exec("DROP TABLE USER(Name, text);")
//	//result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "rou", "si")
//	//if err == nil {
//	//	affected, _ := result.RowsAffected()
//	//	log.Println(affected)
//	//}
//
//	row := db.QueryRow("select Name from User") // 返回一条记录
//	log.Println(row)
//	var name string
//	if err := row.Scan(&name); err == nil {
//		log.Println(name)
//	}
//
//}

//endregion

func main() {
	engine, _ := geeorm.NewEngine("sqlite3", "D:\\Common\\Sqlite\\db\\gee.db")
	defer engine.Close()
	s := engine.NewSession()
	result, _ := s.Raw("insert into User values (?, ?)", "hh", 28).Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
