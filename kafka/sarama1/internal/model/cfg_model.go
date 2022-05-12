package model

//ApplicationConfig : ApplicationConfig
type ApplicationConfig struct {
	ServerConfig ServerConfig `yaml:"server"`
}

//ServerConfig : ServerConfig
type ServerConfig struct {
	Port         string `yaml:"port,omitempty"`
	TimeZone     string `yaml:"time_zone"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
	Loglevel     string `yaml:"loglevel,omitempty"`
	Environment  string `yaml:"env,omitempty"`
	BaseURL      string `yaml:"base_url"`
}
