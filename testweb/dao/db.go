package dao

import (
	"fmt"
	"time"

	"github.com/jiyuwu/gotest/testweb/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type DBConfig struct {
	Port     int      `yaml:"port" validate:"required"`
	Username string   `yaml:"username" validate:"required"`
	Password string   `yaml:"password" validate:"required"`
	DBName   string   `yaml:"dbname" validate:"required"`
	MaxIdle  int      `yaml:"max_idle"`
	MaxOpen  int      `yaml:"max_open"`
	Debug    bool     `yaml:"debug"`
	Master   []string `yaml:"master"` // 主库host
	Slaver   []string `yaml:"slaver"` // 从库host
}

var (
	// 可以根据数据库名，自己重命名此变量
	databaseConn *gorm.DB
	databaseName = "jiyuwuchat"
)

// InitDatabaseConn 这里可以初始化N个数据库链接
func InitDatabaseConn() {
	//配置host 127.0.0.1 mysql.default.svc.cluster.local
	accountConf := &DBConfig{Port: 3306, Username: "ych", Password: "123456", DBName: databaseName, MaxIdle: 4, MaxOpen: 4, Debug: false, Master: []string{"mysql.default.svc.cluster.local"}, Slaver: []string{"mysql.default.svc.cluster.local"}}
	databaseConn = NewClient(accountConf)

	autoMigrate()
}

func autoMigrate() {
	databaseConn.AutoMigrate(new(entity.AccountEntity))
}
func NewClient(c *DBConfig) *gorm.DB {
	if len(c.Master) == 0 {
		log.Fatalf("master db required")
	}

	config := &gorm.Config{
		AllowGlobalUpdate: false,
		NowFunc: func() time.Time {
			return time.Now().UTC().Truncate(time.Millisecond)
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 主库
	master := make([]gorm.Dialector, 0)
	for _, host := range c.Master {
		s := mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
				c.Username, c.Password, host, c.Port, c.DBName),
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		})
		master = append(master, s)
	}

	// 从库，除了host与master不一样，其他配置全部相同
	slaver := make([]gorm.Dialector, 0)
	for _, host := range c.Slaver {
		s := mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
				c.Username, c.Password, host, c.Port, c.DBName),
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		})
		slaver = append(slaver, s)
	}

	db, err := gorm.Open(master[0], config)
	if err != nil {
		log.Fatalf("initDB args:%+v err:%s", c, err)
	}

	if c.Debug {
		db = db.Debug()
	}
	resolver := dbresolver.Register(dbresolver.Config{
		Sources:  master, // 一主
		Replicas: slaver, // 多从
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(10 * time.Minute).
		SetMaxIdleConns(c.MaxIdle).SetMaxOpenConns(c.MaxOpen)
	if err = db.Use(resolver); err != nil {
		log.Fatalf("Use args:%+v err:%s", c, err)
	}
	return db
}
