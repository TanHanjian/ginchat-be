package mydb

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitMySql() {
	new_logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db1, err := gorm.Open(mysql.Open(viper.GetString("mydb.dsn")), &gorm.Config{
		Logger: new_logger,
	})
	DB = db1
	if err != nil {
		fmt.Println(err)
	}
}
