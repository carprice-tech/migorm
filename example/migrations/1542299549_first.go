package migrations

import (
	"errors"

	"github.com/carprice-tech/migorm"
	"github.com/jinzhu/gorm"
)

func init() {
	migorm.RegisterMigration(&migrationFirst{})
}

type migrationFirst struct{}

func (m *migrationFirst) Up(db *gorm.DB, di migorm.MigraterDI) error {
	err := errors.New("implement me")

	return err
}

func (m *migrationFirst) Down(db *gorm.DB, di migorm.MigraterDI) error {
	err := errors.New("implement me")

	return err
}
