package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var ConfInstance Conf

type Conf struct {
	DataSource    DataSource `json:"dataSource"`
}

type DataSource struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Database string `json:"database"`
	LogLevel string `json:"logLevel"`
}

var Ds *gorm.DB

var LoggerLevelMap = map[string]logger.LogLevel{
	"Info":    logger.Info,
	"Error":   logger.Error,
	"Warn":    logger.Warn,
	"Disable": logger.Silent,
}

func InitDb() {
	fmt.Println(ConfInstance)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", ConfInstance.DataSource.User, ConfInstance.DataSource.Password, ConfInstance.DataSource.Host, ConfInstance.DataSource.Database)
	fmt.Println("dns:", dsn)

	var err error
	gormConfig := &gorm.Config{}
	if val, ok := LoggerLevelMap[ConfInstance.DataSource.LogLevel]; ok {
		fmt.Println("设置日志：", ConfInstance.DataSource.LogLevel)
		gormConfig.Logger = logger.Default.LogMode(val)
	}

	Ds, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), gormConfig)
	if err != nil {
		fmt.Println("数据库连接失败:", err.Error())
		return
	}

	sqlDB, _ := Ds.DB()
	sqlDB.SetMaxIdleConns(10)           // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetConnMaxLifetime(time.Hour) // SetConnMaxLifetime 设置了连接可复用的最大时间。
	fmt.Println("mysql init success")
	return
}

