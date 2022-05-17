package config

type Config struct {
	Logging     LoggingConfig `mapstructure:"logging"`
	AWS         AWSConfig     `mapstructure:"aws"`
	Server      ServerConfig  `mapstructure:"server"`
	Backup      BackupConfig  `mapstructure:"backup"`
	DBRegistry  DBRegistry    `mapstructure:"dbregistry,squash"`
	DatabaseURL string        `mapstructure:"databaseURL""`
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
	AccountId            string    `mapstructure:"accountId"`
	RDS                  RDSConfig `mapstructure:"rds"`
}

type RDSConfig struct {
	MasterUsername            string   `mapstructure:"masterUsername"`
	MasterPassword            string   `mapstructure:"masterPassword"`
	SupportedInstanceClasses  []string `mapstructure:"supportedInstanceClasses"`
	SupportedEngines          []string `mapstructure:"supportedEngines"`
	SubnetGroupNames          []string `mapstructure:"subnetGroupNames"`
	AllowedStorageSizeGBRange string   `mapstructure:"allowedStorageSizeGBRange"`
}

type LoggingConfig struct {
	Level   string `mapstructure:"level"`
	JSONLog bool   `mapstructure:"jsonLogs"`
}
