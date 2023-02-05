package conf

import (
	"fmt"
	"io/ioutil"

	"experiment_go/kafka/dyn_consumer/internal/data/model"

	yaml "gopkg.in/yaml.v2"
)

func Init(path string) (*model.AppConfig, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(model.AppConfig)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	if err := setEnvironmentFromAppConfig(cfg); err != nil {
		return nil, fmt.Errorf("unable to set environment. err: %v", err)
	}

	if err := setMerchantIdFromAppConfig(cfg); err != nil {
		return nil, fmt.Errorf("unable to set merchant id. err: %v", err)
	}

	return cfg, nil
}
