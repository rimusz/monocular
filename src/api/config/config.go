package config

import (
	"os"

	"github.com/imdario/mergo"
)

// ConfigurationWithOverrides includes default Configuration values
// and environment specific ones
type configurationWithOverrides map[string]Configuration

// Configuration is the the resulting environment based Configuration
// For now it only includes Cors info
type Configuration struct {
	Cors
}

// Cors configuration used during middleware setup
type Cors struct {
	AllowedOrigins []string
	AllowedHeaders []string
}

func currentEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "production"
	}
	return env
}

// The default configuration gets overridden by the `currentEnvironment` settings (if any)
func readConfigWithOverrides() configurationWithOverrides {
	var config = configurationWithOverrides{
		"default": Configuration{
			Cors{
				AllowedOrigins: []string{"my-api-server"},
				AllowedHeaders: []string{"access-control-allow-headers", "x-xsrf-token"},
			},
		},
		"development": Configuration{
			Cors{
				AllowedOrigins: []string{"*"},
			},
		},
	}

	return config
}

// GetConfig returns the environment specific configuration
func GetConfig() (Configuration, error) {
	config := configurationWithOverrides{}
	res := Configuration{}

	config = readConfigWithOverrides()

	res = mergeConfig(config, currentEnvironment())

	return res, nil
}

func mergeConfig(conf configurationWithOverrides, env string) Configuration {
	defaults := conf["default"]
	custom := conf[env]

	mergo.Merge(&custom, defaults)
	return custom
}
