package database

import (
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

type Database struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (d *Database) GetRandomNumber() (number int, err error) {
	d.logger.Debug("Database.GetRandomNumber")

	var result struct {
		Number int
	}

	// Do some calculations to consume resources
	res := d.db.Raw(`
WITH RECURSIVE foo(cur, nex) AS ( 
	SELECT 1,1 
		UNION ALL 
	SELECT nex, cur+1 
	FROM foo 
	LIMIT 100000
) 	SELECT sum(cur) AS number
	FROM foo;
`).Scan(&result)
	if res.Error != nil {
		err = fmt.Errorf("res.Error: %w", res.Error)
		return
	}

	number = result.Number
	return
}

func New(logger *zap.Logger) (database *Database, err error) {
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault()
	gormLogger.LogLevel = gormlogger.Warn

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = path.Join(os.TempDir(), "go-fx-test.db")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		err = fmt.Errorf("gorm.Open: %w", err)
		return
	}

	database = new(Database)
	database.db = db
	database.logger = logger

	return
}
