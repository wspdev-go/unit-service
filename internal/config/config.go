package config

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
	Host     string `yaml:"Host" json:"host"`
	Port     int    `yaml:"Port" json:"port"`
	Username string `yaml:"Username" json:"username"`
	Password string `yaml:"Password" json:"password"`
	Database string `yaml:"Database" json:"database"`
}

type TransactionConfig struct {
	Host     string `yaml:"Host" json:"host"`
	Port     int    `yaml:"Port" json:"port"`
	Username string `yaml:"Username" json:"username"`
	Password string `yaml:"Password" json:"password"`
	Database string `yaml:"Database" json:"database"`
}
