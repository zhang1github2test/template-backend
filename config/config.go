package config

import (
	"fmt"
	"template-backend/internal/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
)

type AppConfig struct {
	App struct {
		Name string
		Port int
		Env  string
	} `mapstructure:"app"`

	Database struct {
		Host     string
		DbName   string
		Account  string
		Password string
	} `mapstructure:"database"`

	JWT struct {
		Secret       string
		Expires      int
		SkipAuthUrls []string `mapstructure:"skip_auth_urls"`
	} `mapstructure:"jwt"`
}

var (
	cfg *AppConfig
)

// LoadConfig 从指定目录加载配置，如果没传目录则默认当前目录
func LoadConfig(path ...string) *AppConfig {
	viper.SetConfigName("config") // 配置文件名: config.yaml
	viper.SetConfigType("yaml")

	// 如果传了 path[0]，就用传入的目录，否则用当前目录
	if len(path) > 0 {
		viper.AddConfigPath(path[0])
	} else {
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // 支持环境变量覆盖

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	c := &AppConfig{}
	if err := viper.Unmarshal(c); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	fmt.Println("配置加载成功:", viper.ConfigFileUsed())
	cfg = c
	return cfg
}

func GetConfig() *AppConfig {
	if cfg == nil {
		log.Fatal("配置未初始化，请先调用 LoadConfig()")
	}
	return cfg
}

func InitDB() *gorm.DB {
	database := GetConfig().Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", database.Account, database.Password, database.Host, database.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	// 自动建表
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Role{})
	db.AutoMigrate(&model.Menu{})
	db.AutoMigrate(&model.Config{})
	db.AutoMigrate(&model.Log{})
	return db
}
