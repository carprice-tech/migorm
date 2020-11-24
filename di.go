package migorm

func NewMigraterDI(log Logger) MigraterDI {
	return &migraterDI{
		Log: log,
	}
}

type MigraterDI interface {
	GetLogger() Logger
}
type migraterDI struct {
	Log Logger
}

func (di *migraterDI) GetLogger() Logger {
	return di.Log
}
