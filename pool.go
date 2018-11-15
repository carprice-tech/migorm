package migorm

import (
	"github.com/jinzhu/gorm"
	"runtime"
	"strings"
	"path/filepath"
	"sync"
)

var pool migrationsPool

type migrationsPool struct {
	migrations map[string]Migration
	sync.Mutex
}

func init() {
	pool = migrationsPool{migrations: make(map[string]Migration)}
}

type Migration interface {
	Up(db *gorm.DB, log Logger) error
	Down(db *gorm.DB, log Logger) error
}

// Each migration file call this method in its init method
func RegisterMigration(migration Migration) {

	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Fail invoke caller")
		return
	}
	migrationName := strings.Replace(filepath.Base(file), ".go", "", -1)

	pool.Lock()
	defer pool.Unlock()
	_, ok = pool.migrations[migrationName]
	if (ok) {
		panic("Migration with name : " + migrationName + " already exist")
	}
	pool.migrations[migrationName] = migration
}