package config

type Config struct {
	Logging    LoggingConfig `mapstructure:"logging"`
	AWS        AWSConfig     `mapstructure:"aws"`
	Server     ServerConfig  `mapstructure:"server"`
	Backup     BackupConfig  `mapstructure:"backup"`
	DBRegistry DBRegistry    `mapstructure:"dbregistry,squash"`
}

type BackupConfig struct {
	S3Bucket string `mapstructure:"s3_bucket"`
}

type DBRegistry struct {
	DBs []DB `mapstructure:"databases"`
}

type DB struct {
	Host       string `mapstructure:"host"`
	Port       int32  `mapstructure:"port"`
	DBName     string `mapstructure:"dbName"`
	DBUser     string `mapstructure:"user"`
	DBPassword string `mapstructure:"password"`
}

type ServerConfig struct {
	BindAddress string `mapstructure:"bindAddress"`
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
