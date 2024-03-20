package config

import (
	"os"

	"github.com/joho/godotenv"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	}
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"-"`
		Dbname   string `yaml:"dbname"`
	}
}

const (
	configFile          = "config.yml"
	databasePasswordEnv = "DATABASE_PASSWORD"
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
	return cfg, nil
}
