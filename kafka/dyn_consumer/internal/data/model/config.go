package model

type AppConfig struct {
	Server ServerConfig `yaml:"server"`
	Logger LogConfig    `yaml:"logger"`
}

type ServerConfig struct {
	Port         string `yaml:"port,omitempty"`
	TimeZone     string `yaml:"time_zone"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
	Loglevel     string `yaml:"loglevel,omitempty"`
	BaseURL      string `yaml:"base_url"`
	Environment  string `yaml:"env,omitempty"`
	MerchantID   string `yaml:"merchant_id"`
}

type LogConfig struct {
	EnableStdout   bool            `yaml:"enable_stdout"`
	EnableLogFile  bool            `yaml:"enable_logfile"`
	CallerSkipSet  bool            `yaml:"caller_skipset"`
	CallerSkip     int             `yaml:"caller_skip"`
	LogFileConfigs []LogFileConfig `yaml:"logfile_configs"`
}

type LogFileConfig struct {
	Levels           []string `yaml:"levels"`
	IsAccessLog      bool     `yaml:"is_access_log"`
	FullpathFilename string   `yaml:"fullpath_filename"`
	MaxSize          int      `yaml:"max_size"`
	MaxAge           int      `yaml:"max_age"`
	MaxBackups       int      `yaml:"max_backups"`
	LocalTime        bool     `yaml:"local_time"`
	Compress         bool     `yaml:"compress"`
}
