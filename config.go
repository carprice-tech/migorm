package migorm

type Configurator struct {
	Log           Logger
	MigrationsDir string
	TableName     string
}