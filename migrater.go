package migorm

import (
	"github.com/jinzhu/gorm"
	"sort"
	"errors"
	"time"
	"strings"
	"io/ioutil"
	"os"
	"fmt"
	"text/template"
	"runtime"
	"path"
)

func NewMigrater(db *gorm.DB) Migrater {
	return &migrater{
		db: db,
		Configurator: &Configurator{
			Log:           NewLogger(),
			MigrationsDir: "migrations",
			TableName:     "migrations",
		},
	}
}

type Migrater interface {
	Conf() *Configurator
	UpMigrations() error
	UpConcreteMigration(name string) error
	DownConcreteMigration(name string) error
	MakeFileMigration(name string) error
}

type migrater struct {
	db *gorm.DB
	*Configurator
}

func (m *migrater) Conf() *Configurator {
	return m.Configurator
}

func (m *migrater) UpMigrations() error {

	m.Log.Infof("Start migrations")

	m.checkMigrationTable()

	newMigrations := m.getNewMigrations()

	successCnt := 0
	for _, migration := range newMigrations {
		if migration.Id == 0 {
			tx := m.db.Begin()
			if err := pool.migrations[migration.Name].Up(tx, m.Log); err != nil {
				tx.Rollback()
				return fmt.Errorf("up migration: %+v, err: %+v", migration.Name, err)
			}
			if err := m.db.Create(&migration).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("save migration: %v, err: %+v", migration.Name, err)
			}
			tx.Commit()
			m.Log.Infof("success: %+v", migration.Name)
			successCnt++
		}
	}

	if successCnt > 0 {
		m.Log.Infof("All migrations are done success!")
	} else {
		m.Log.Infof("Nothing to migrate.")
	}

	return nil
}

func (m *migrater) UpConcreteMigration(name string) error {
	mig, ok := pool.migrations[name]
	if !ok {
		return errors.New("Does not exist migration with name: " + name)
	}

	tx := m.db.Begin()
	if err := mig.Up(tx, m.Log); err != nil {
		tx.Rollback()
		return err
	}

	migrationModel := m.newMigrationModel()
	err := m.db.Where("name = ?", name).First(&migrationModel).Error
	if !gorm.IsRecordNotFoundError(err) && err != nil {
		return err
	}

	if migrationModel.Id == 0 {
		migrationModel.Name = name
		if err := m.db.Create(&migrationModel).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (m *migrater) DownConcreteMigration(name string) error {

	mig, ok := pool.migrations[name]
	if !ok {
		return errors.New("Does not exist migration with name: " + name)
	}

	tx := m.db.Begin()
	if err := mig.Down(tx, m.Log); err != nil {
		tx.Rollback()
		return err
	}

	migrationModel := m.newMigrationModel()
	err := m.db.Where("name = ?", name).First(&migrationModel).Error
	if !gorm.IsRecordNotFoundError(err) && err != nil {
		return err
	}

	if migrationModel.Id != 0 {
		if err := m.db.Delete(&migrationModel, "name = ?", name).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (m *migrater) MakeFileMigration(name string) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsPath := currentDir + "/" + m.Configurator.MigrationsDir

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		m.Log.Infof("Create new directory : %v", migrationsPath)
		if err := os.MkdirAll(migrationsPath, os.ModePerm); err != nil {
			return err
		}
	}

	err = checkFileExists(migrationsPath, name+".go")
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	realName := fmt.Sprintf("%d_%s.go", now, name)

	migrationPath := migrationsPath + "/" + realName

	f, err := os.Create(migrationPath)
	if err != nil {
		return fmt.Errorf("create migration file: %v", err)
	}

	partsName := strings.Split(name, "_")
	structName := "migration"
	for _, p := range partsName {
		structName += strings.Title(p)
	}

	partsDir := strings.Split(m.Configurator.MigrationsDir, "/")
	packageName := partsDir[len(partsDir)-1]

	tmpl, err := getTemplate()
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, map[string]interface{}{"struct_name": structName, "package": packageName})

	if err != nil {
		return err
	}

	m.Log.Infof("migration file created: %v", realName)

	return nil
}

// Finds not yet completed migration files
func (m *migrater) getNewMigrations() []migrationModel {

	var names []string
	for k, _ := range pool.migrations {
		names = append(names, k)
	}

	sort.Strings(names)

	step := 20 // limit
	result := make([]migrationModel, 0)
	for i := 0; i < len(names); {

		i += step
		var chunkNames []string
		if i <= len(names) {
			chunkNames = names[i-step : i]
		} else {
			chunkNames = names[i-step:]
		}

		rows := make([]struct{ Name string }, 0)
		if err := m.db.Model(m.newMigrationModel()).
			Where("name IN (?)", chunkNames).
			Scan(&rows).Error; err != nil {

			panic(err)
		}
		existMigrations := make(map[string]bool)
		for _, row := range rows {
			existMigrations[row.Name] = true
		}

		for _, name := range names {
			if _, ok := existMigrations[name]; !ok {
				model := m.newMigrationModel()
				model.Name = name
				result = append(result, model)
			}
		}
	}

	return result
}

//
func (m *migrater) newMigrationModel() migrationModel{
	return migrationModel{tableName: m.Configurator.TableName}
}

// ***  helpers ***

// check or create table to register successful migrations
func (m *migrater) checkMigrationTable() {

	model := m.newMigrationModel()

	if !m.db.HasTable(&model) {
		m.Log.Infof("Init table: %v", model.TableName())
		if err := m.db.AutoMigrate(&model).Error; err != nil {
			panic(err)
		}
	}
}

// сheck the existence of a file in the directory with migrations
func checkFileExists(dir string, name string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		split := strings.Split(f.Name(), "_")

		if name == strings.Join(split[1:], "_") {
			return fmt.Errorf("File %v already exists in dir: %v", name, dir)
		}
	}

	return nil
}

//
func getTemplate() (*template.Template, error) {

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("Template caller")
	}

	tmpl, err := template.ParseFiles(path.Dir(filename) + "/" + "template")
	if err != nil {
		return nil, fmt.Errorf("parse template : %v", err)
	}

	return tmpl, nil
}