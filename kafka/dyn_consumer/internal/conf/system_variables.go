package conf

import (
	"fmt"
	"strings"

	"experiment_go/kafka/dyn_consumer/internal/data/model"
	"experiment_go/kafka/dyn_consumer/internal/pkg/constant"
)

var (
	env        string = constant.SystemValueCodeEnvDev
	merchantId string = ""
)

func MerchantID() string {
	return merchantId
}

func ENV() string {
	return env
}

func setMerchantIdFromAppConfig(cfg *model.AppConfig) error {
	if cfg == nil {
		return fmt.Errorf("supplied AppConfig is nil")
	}

	if cfg.Server.MerchantID == "" {
		return fmt.Errorf("merchant id should not be empty")
	}

	merchantId = cfg.Server.MerchantID

	return nil
}

func setEnvironmentFromAppConfig(cfg *model.AppConfig) error {
	if cfg == nil {
		return fmt.Errorf("supplied AppConfig is nil")
	}

	environment := cfg.Server.Environment
	environment = strings.ToUpper(environment)

	if environment == "PRODUCTION" {
		environment = constant.SystemValueCodeEnvProd
	}

	switch environment {
	case constant.SystemValueCodeEnvDev,
		constant.SystemValueCodeEnvTest,
		constant.SystemValueCodeEnvUAT,
		constant.SystemValueCodeEnvStaging,
		constant.SystemValueCodeEnvProd,
		constant.SystemValueCodeEnvDevProd:

		env = environment

		return nil
	}

	return fmt.Errorf("environment [ %v ] is not a valid environment", cfg.Server.Environment)
}
