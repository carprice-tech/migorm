package main

import (
	"github.com/carprice-tech/migorm"
	_ "github.com/carprice-tech/migorm/example/migrations"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"path"
	"runtime"
)

//
// only for run  make migrations command
// uncomment lines for create dbConn and fill connection params for execute other commands
func main() {
	var dbConn *gorm.DB

	// db_user := ""
	// db_pass := ""
	// db_host := ""
	// db_port := ""
	// db_name := ""
	//
	// conStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	// dbConn, err := gorm.Open("mysql", conStr)
	//
	// if err != nil{
	//   panic(err)
	// }

	migrater := migorm.NewMigrater(
		dbConn,
		migorm.NewMigraterDI(migorm.NewLogger()),
	)

	_, file, _, _  := runtime.Caller(0)
	curDir := path.Dir(file)
	migrater.SetMigrationsDir(curDir + "/../migrations")

	migorm.Run(migrater)
}
