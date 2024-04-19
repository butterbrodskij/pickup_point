package config

import (
	"os"

	"github.com/joho/godotenv"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port       string `yaml:"port"`
		SecurePort string `yaml:"secure_port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"-"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`
	Users []struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"users"`
	Kafka struct {
		Brokers []string `yaml:"brokers,omitempty"`
		Topic   string   `yaml:"topic"`
	} `yaml:"kafka"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DB       int    `yaml:"db"`
		Password string `yaml:"-"`
	} `yaml:"redis"`
}

const (
	configFile          = "config.yml"
	databasePasswordEnv = "DATABASE_PASSWORD"
	redisPasswordEnv    = "REDIS_PASSWORD"
	QueryParamKey       = "point"
	CertFile            = "server.crt"
	KeyFile             = "server.key"
)

func GetConfig() (Config, error) {
	var cfg Config
	rawBytes, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(rawBytes, &cfg)
	if err != nil {
		return Config{}, err
	}
	if err = godotenv.Load(); err != nil {
		return Config{}, err
	}
	pass, ok := os.LookupEnv(databasePasswordEnv)
	if !ok {
		return Config{}, model.ErrorInvalidEnvironment
	}
	cfg.Database.Password = pass
	passRedis, ok := os.LookupEnv(redisPasswordEnv)
	if !ok {
		return Config{}, model.ErrorInvalidEnvironment
	}
	cfg.Redis.Password = passRedis
	return cfg, nil
}
