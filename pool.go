package migorm

import (
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

var pool migrationsPool // nolint:gochecknoglobals

type migrationsPool struct {
	migrations map[string]Migration
	sync.Mutex
}

func init() {
	pool = migrationsPool{migrations: make(map[string]Migration)}
}

type Migration interface {
	Up(db *gorm.DB, di MigraterDI) error
	Down(db *gorm.DB, di MigraterDI) error
}

// Each migration file call this method in its init method
func RegisterMigration(migration Migration) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Fail invoke caller")
	}
	migrationName := strings.Replace(filepath.Base(file), ".go", "", -1) // nolint:gocritic

	pool.Lock()
	defer pool.Unlock()
	_, ok = pool.migrations[migrationName]
	if ok {
		panic("Migration with name : " + migrationName + " already exist")
	}
	pool.migrations[migrationName] = migration
}
