package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Database struct {
	dbs     map[string]*gorm.DB
	configs map[string]DatabaseConfig
	logger  *log.Logger
}

var opens = map[string]func(string) gorm.Dialector{
	"mysql":    mysql.Open,
	"sqlite":   sqlite.Open,
	"postgres": postgres.Open,
}

type DatabaseConfig struct {
	Name   string
	Dsn    string
	Driver string
	Tables []Table
}

func NewDatabase(configs []DatabaseConfig) *Database {
	db := &Database{
		dbs:     make(map[string]*gorm.DB),
		configs: make(map[string]DatabaseConfig),
		logger:  log.New(os.Stdout, "|Database| ", log.LstdFlags),
	}

	for _, config := range configs {
		db.AddDatabase(config)
	}

	return db
}

// AddDatabase 增加数据库配置
func (db *Database) AddDatabase(config DatabaseConfig) {
	// db.configs = append(db.configs, config)
	db.configs[config.Name] = config
}

// Init 初始化数据库连接
func (db *Database) Init() error {
	if len(db.configs) == 0 {
		return fmt.Errorf("no database config")
	}

	var newLogger glogger.Interface
	newLogger = CustomLogger{
		glogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			glogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  glogger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	for _, config := range db.configs {
		db.logger.Printf("open database: %s", config.Name)
		conn, err := gorm.Open(opens[config.Driver](config.Dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			return fmt.Errorf("open database %s error: %v", config.Name, err)
		}
		db.dbs[config.Name] = conn

		// migrate tables
		for _, table := range config.Tables {
			db.logger.Printf("migrate table: %s", table.Name)
			err = conn.AutoMigrate(table.Definition)
			if err != nil {
				return fmt.Errorf("migrate table %s error: %v", table.Name, err)
			}
		}
	}

	return nil
}

// GetDatabase 获取数据库连接
func (db *Database) GetDatabase(name string) *gorm.DB {
	return db.dbs[name]
}

type Table struct {
	Name       string
	Definition any
}

type CustomLogger struct {
	glogger.Interface
}

func (c CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[INFO] "+msg, data...)
}

func (c CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[WARN] "+msg, data...)
}

func (c CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[ERROR] "+msg, data...)
}

func (c CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 对于复杂的查询参数，你可能需要进一步处理
	// 这里假设所有参数都是简单的字符串、数字等
	for _, param := range glogger.ExplainSQL(sql, nil, "") {
		sql = strings.Replace(sql, "?", fmt.Sprintf("'%v'", param), 1)
	}

	if err != nil {
		log.Printf("[rows:%v] SQL: %s |  Elapsed: %v | Error: %v", rows, sql, elapsed, err)
	} else {
		log.Printf("[rows:%v] SQL: %s |  Elapsed: %v ", rows, sql, elapsed)
	}
}
