package utils

import (
	"fmt"
	"ginchat/mydb"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"regexp"
)

var Go_validate *validator.Validate

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config app:", viper.Get("mydb"))
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

func Init() {
	InitConfig()
	mydb.InitMySql()
	InitValidator()
}
