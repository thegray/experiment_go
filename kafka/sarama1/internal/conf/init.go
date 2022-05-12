package conf

import (
	"experiment_go/kafka/sarama1/internal/model"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var conf = new(model.ApplicationConfig)

//InitServiceConfig : Init service config from local file path
func InitServiceConfig(appName string) error {
	configPath := genServiceConfigPath(appName)
	err := parseConfigFromPath(configPath)
	if err != nil {
		return err
	}
	return nil
}

func genServiceConfigPath(appName string) string {
	// region = env.GetRegion().String()
	// envv := env.GetEnv()
	// return fmt.Sprintf("./conf/%s/%s/id-maker.%s.toml", appName, region, envv)
	return fmt.Sprintf("./conf/%s/config.yaml", appName)
}

func parseConfigFromPath(configPath string) error {
	// fmt.Printf("get service config from config path: %v\n", configPath)
	// if _, err := toml.DecodeFile(configPath, conf); err != nil {
	// 	fmt.Printf("failed to decode config file, error - %s", err.Error())
	// 	return err
	// }
	// fmt.Printf("get service config:%+v\n", conf)
	// return nil
	fmt.Printf("get service config from config path: %v\n", configPath)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}
	if err := yaml.Unmarshal(bytes, conf); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	return nil
}

//GetGlobalConfig : Get global service config from local
func GetGlobalConfig() *model.ApplicationConfig {
	return conf
}

//SetGlobalConfig : Set global service config, only test
func SetGlobalConfig(c *model.ApplicationConfig) {
	conf = c
}
