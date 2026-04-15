package config

import (
	"os"
	"unit-service/logger"

	"gopkg.in/yaml.v3"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ReferenceDB   *ReferenceConfig   `yaml:"reference-db"`
	QueueDB       *QueueConfig       `yaml:"queue-db"`
	TransactionDB *TransactionConfig `yaml:"transaction-db"`
}

type ReferenceConfig struct {
	Host     string `yaml:"Host" json:"host"`
	Port     int    `yaml:"Port" json:"port"`
	Username string `yaml:"Username" json:"username"`
	Password string `yaml:"Password" json:"password"`
	Database string `yaml:"Database" json:"database"`
}

type QueueConfig struct {
	Host       string   `yaml:"Host" json:"host"`
	Port       int      `yaml:"Port" json:"port"`
	Username   string   `yaml:"Username" json:"username"`
	Password   string   `yaml:"Password" json:"password"`
	Database   int      `yaml:"Database" json:"database"`
	PoolSize   int      `yaml:"PoolSize" json:"poolSize"`
	MasterName string   `yaml:"MasterName,omitempty" json:"masterName,omitempty"`
	Addresses  []string `yaml:"Addresses,omitempty" json:"addresses,omitempty"`
}

type TransactionConfig struct {
	Host         string `yaml:"Host" json:"host"`
	Port         int    `yaml:"Port" json:"port"`
	Username     string `yaml:"Username" json:"username"`
	Password     string `yaml:"Password" json:"password"`
	Database     string `yaml:"Database" json:"database"`
	DialTimeout  int    `yaml:"DialTimeout" json:"dialTimeout"`
	MaxOpenConns int    `yaml:"MaxOpenConns" json:"maxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns" json:"maxIdleConns"`
}

func GetConfig(configPath string) (*Config, error) {
	var cfg Config

	err := cfg.readFile(configPath)
	if err != nil {
		return nil, err
	}

	err = cfg.readEnv()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (r *Config) readFile(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		logger.Error("Can't open config file: %s", configPath)
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			logger.Error("Can't close config file: %s", err)
		}
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(r)
	if err != nil {
		logger.Error("Can't decode config file: %s", configPath)
		return err
	}

	return nil
}

func (r *Config) readEnv() error {
	err := envconfig.Process("", r)
	if err != nil {
		logger.Error("Can't read environment variables: %s", err)
	}
	return err
}
