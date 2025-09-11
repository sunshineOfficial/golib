package db

import (
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"

	"github.com/pressly/goose/v3"
	"github.com/sunshineOfficial/golib/golog"
)

type gooseLogger struct {
	log golog.Logger
}

func (s gooseLogger) Fatal(v ...interface{}) {
	s.log.ErrorEntry(fmt.Sprint(v...)).WithTags("migration").Write()
}

func (s gooseLogger) Fatalf(format string, v ...interface{}) {
	s.log.ErrorEntryf(format, v...).WithTags("migration").Write()
}

func (s gooseLogger) Print(v ...interface{}) {
	s.log.Debug(fmt.Sprint(v...))
}

func (s gooseLogger) Printf(format string, v ...interface{}) {
	s.log.Debugf(format, v...)
}

func (s gooseLogger) Println(v ...interface{}) {
	s.log.Debug(fmt.Sprintln(v...))
}

// Migrate применяет доступные миграции из файлов по пути path в указанной fs
func Migrate(fs fs.FS, log golog.Logger, db *sqlx.DB, path string) error {
	goose.SetVerbose(false)
	goose.SetLogger(gooseLogger{log: log})
	goose.SetBaseFS(fs)
	goose.SetTableName("db_version")

	return goose.Up(db.DB, path)
}
