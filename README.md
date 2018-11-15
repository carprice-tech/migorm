# MIGORM

MIGORM - Helper tool for [gorm](https://github.com/jinzhu/gorm) framework, allows you to make changes in your database,
 by creating migrations files.

## Quick start:
```bash
glide get github.com/carprice-tech/migorm
```
Then you need to create a file for run migrations commands (e.g. migrate.go)
```go
package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/carprice-tech/migorm"
)

func main() {

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
```

For testing, try create a new migration.

```bash
go run migrate.go make first_example
```

After that, a package with default name: **migrations**,  will be created in the same directory. And the migration file (**<timestamp>_first_example.go**) will be created there.

```
.
├── cmd
│   ├── migrate.go
│   └── migrations
│       └── 1542223496_first_example.go
...
```

**! You must import new created package** with migrations in you migration run file (migrate.go)
```go

import (
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/jinzhu/gorm"
    "github.com/carprice-tech/migorm"
    _ "your/project/path/migrations"
)
```

Now all you need is to insert your code into the up/down methods in new created file (<timestamp>_first_example.go)


Run migrations for execution.

```bash
go run migrate.go
```
Successful migrations are inserted into the table by default: **migrations**. And in the future will be skipped

## Available commands:
```bash
go run migrate.go                                # run new migrations
go run migrate.go make my_new_migration          # create migration
go run migrate.go up 1542223496_first_example    # Up specific migration
go run migrate.go down 1542223496_first_example  # Down specific migration
```


## Configuration
You can configure some parameters before performing migrations.
```go

migrater := migorm.NewMigrater(db)

// don't forget import new package after create
migrater.Conf().MigrationsDir = "../../my_migration_pkg" // relative current file

migrater.Conf().TableName = "my_migrations"

migrater.Conf().Log = migorm.NewLogger() // or your implementation

migorm.Run(migrater)

```


## Recomendations
// TODO

## Authors

+ [pzavyalov](https://github.com/pzavyalov)

## License
This project is licensed under the MIT License