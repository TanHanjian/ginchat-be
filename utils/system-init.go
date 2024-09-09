package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var Go_validate *validator.Validate

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config app:", viper.Get("mysql"))
}

func InitMySql() {
	new_logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db1, err := gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{
		Logger: new_logger,
	})
	DB = db1
	if err != nil {
		fmt.Println(err)
	}
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 正则表达式匹配中国手机号码
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}

func InitValidator() {
	Go_validate = validator.New(validator.WithRequiredStructEnabled())
	Go_validate.RegisterValidation("phone", validatePhone)
}
