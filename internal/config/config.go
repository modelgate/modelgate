package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 版本信息，在编译时自动生成
var (
	Version   = "unknown"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type Config struct {
	ApiServer   serverConfig    `envPrefix:"API_SERVER_"`
	AdminServer serverConfig    `envPrefix:"ADMIN_SERVER_"`
	Logger      loggerConfig    `envPrefix:"LOGGER_"`
	Database    databaseConfig  `envPrefix:"DATABASE_"`
	JWT         JWT             `envPrefix:"JWT_"`
	Temp        tempConfig      `envPrefix:"TEMP_"`
	Secret      secretConfig    `envPrefix:"SECRET_"`
	RateLimit   RateLimitConfig `envPrefix:"RATE_LIMIT_"`
	Redis       redisConfig     `envPrefix:"REDIS_"`
}

type databaseConfig struct {
	Type     string `env:"TYPE"`
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

type redisConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}

type secretConfig struct {
	Key string `env:"KEY"`
}

type tempConfig struct {
	Dir string `env:"DIR"`
}

type loggerConfig struct {
	Format string
	Level  string
	Output string
}

type JWT struct {
	Key string
}

type serverConfig struct {
	Name  string
	Host  string
	Port  int
	Debug bool
}

type RateLimitConfig struct {
	Enabled           bool    `env:"ENABLED"`
	RequestsPerSecond float64 `env:"REQUESTS_PER_SECOND"`
	Burst             int     `env:"BURST"`
	TrustProxy        bool    `env:"TRUST_PROXY"`
	CleanupInterval   int     `env:"CLEANUP_INTERVAL"`
	ExpireAfter       int     `env:"EXPIRE_AFTER"`
}

var appPath string
var config *Config

// GetConfig 获取配置
func GetConfig() *Config {
	return config
}

// Init 初始化配置
func Init(i do.Injector, _appPath string) {
	var err error
	appPath, err = filepath.Abs(_appPath)
	if err != nil {
		log.Fatalf("init config err: %v", err)
	}

	do.Provide(i, func(i do.Injector) (*Config, error) {
		return NewConfig(filepath.Join(appPath, "configs/config.toml"))
	})

	config := do.MustInvoke[*Config](i)
	// 日志设置
	if err := initLogger(config.Logger); err != nil {
		log.Fatalf("init logger err: %v", err)
	}
	// 注册db
	do.Provide(i, func(i do.Injector) (dbConn *gorm.DB, err error) {
		cfg := config.Database
		switch cfg.Type {
		case "mysql":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?collation=utf8mb4_general_ci&parseTime=True&loc=Local&multiStatements=true&timeout=10s&readTimeout=30s&writeTimeout=30s",
				cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
			dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				return
			}
			return dbConn.Debug(), nil
		default:
			err = errors.New("unsupported database type")
			return
		}
	})
	// 注册redis
	do.Provide(i, func(i do.Injector) (*redis.Client, error) {
		cfg := config.Redis
		client := redis.NewClient(&redis.Options{
			Addr:         cfg.Host + ":" + strconv.Itoa(cfg.Port),
			Password:     cfg.Password,
			DB:           cfg.DB,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
		return client, nil
	})
}

// NewConfig 实例化配置
func NewConfig(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}
	// 将env配置加载到环境变量
	if cErr := gotenv.Load(filepath.Join(appPath, ".env")); cErr != nil {
		log.Debug("load env err:", cErr)
	}
	// 读取环境变量
	if err := env.ParseWithOptions(config, env.Options{Prefix: "MG_", UseFieldNameByDefault: true}); err != nil {
		return nil, err
	}
	if err = initLogger(config.Logger); err != nil {
		return nil, err
	}
	return config, nil
}

func initLogger(cfg loggerConfig) error {
	if cfg.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	lvl, err := log.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	if cfg.Output != "" {
		logFile, err := os.OpenFile(cfg.Output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		log.SetOutput(logFile)
		log.RegisterExitHandler(func() {
			log.Info("退出logger......")
			logFile.Close()
		})
	}
	return nil
}

func GetPath(path string) string {
	return filepath.Join(appPath, path)
}
