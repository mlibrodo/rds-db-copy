package config

type Config struct {
	Logging LoggingConfig `mapstructure:"logging"`
	AWS     AWSConfig     `mapstructure:"aws"`
	Server  ServerConfig  `mapstructure:"server"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type AWSConfig struct {
	ServiceConfigProfile string    `mapstructure:"serviceConfigProfile"`
	RDS                  RDSConfig `mapstructure:"rds"`
}

type RDSConfig struct {
	MasterUsername     string `mapstructure:"masterUsername"`
	MasterUserPassword string `mapstructure:"masterPassword"`
}

type LoggingConfig struct {
	Level   string `mapstructure:"level"`
	JSONLog bool   `mapstructure:"jsonLogs"`
}
