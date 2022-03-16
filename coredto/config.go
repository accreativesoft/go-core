package coredto

type Config struct {
	DatabaseConfig DatabaseConfig `json:"databaseConfig,omitempty"`
	LogConfig      LogConfig      `json:"logConfig,omitempty"`
}

type DatabaseConfig struct {
	Dialect      string `json:"dialect,omitempty"`
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	Database     string `json:"database,omitempty"`
	User         string `json:"user,omitempty"`
	Password     string `json:"password,omitempty"`
	Charset      string `json:"charset,omitempty"`
	LogLevel     string `json:"logLevel,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
	MaxIdleConns int    `json:"maxIdleConns,omitempty"`
	MaxOpenConns int    `json:"maxOpenConns,omitempty"`
}

type LogConfig struct {
	Output     string `json:"output,omitempty"`
	Path       string `json:"path,omitempty"`
	Level      string `json:"level,omitempty"`
	MaxSize    int    `json:"maxSize,omitempty"`
	MaxBackups int    `json:"maxBackups,omitempty"`
	MaxAge     int    `json:"maxAge,omitempty"`
	Compress   bool   `json:"compress,omitempty"`
}
