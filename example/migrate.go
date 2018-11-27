package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/carprice-tech/migorm"
	_ "github.com/carprice-tech/migorm/example/migrations"
)

func main() {

	var dbConn *gorm.DB

	db_user := ""
	db_pass := ""
	db_host := ""
	db_port := ""
	db_name := ""

	conStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	dbConn, err := gorm.Open("mysql", conStr)

	if err != nil{
		panic(err)
	}

	migrater := migorm.NewMigrater(dbConn)
	migorm.Run(migrater)
}