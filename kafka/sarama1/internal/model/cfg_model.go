package model

//ApplicationConfig : ApplicationConfig
type ApplicationConfig struct {
	ServerConfig ServerConfig `yaml:"server"`
	SaramaConfig SaramaConfig `yaml:"sarama"`
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

type SaramaConfig struct {
	BrokersHost []string `yaml:"brokers"`
	Log         bool     `yaml:"log"`
	CertFile    string   `yaml:"certfile"`
	KeyFile     string   `yaml:"keyfile"`
	CaFile      string   `yaml:"cafile"`
	VerifySSL   bool     `yaml:"verifyssl"`
	Topic       string   `yaml:"topic"`
}
