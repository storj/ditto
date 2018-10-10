package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Reads config from default config location($HOME/.ditto/config.json or from user-defined location
// Returns parsed config or error
func ReadDefaultConfig(useDefaults bool) (config *Config, err error) {
	config, err = parseConfig(useDefaults)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	if config.Server1 == nil || config.Server1.IsEmpty() ||
		config.Server2 == nil || config.Server2.IsEmpty() {

		return nil, errors.New("Credentials are not set. Please define credentials with `ditto config set`")
	}

	bytes, _ := json.MarshalIndent(config, "", "\t")
	println(string(bytes))

	return config, nil
}

// Gather config values and applies default values
func parseConfig(useDefaults bool) (config *Config, err error) {
	if viper.IsSet("configPath") {
		viper.SetConfigFile(viper.GetString("configPath"))
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.ditto")
		viper.SetConfigType("json")
	}

	if useDefaults {
		setDefaults()
	}

	// viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func setDefaults() {
	// Root defaults
	viper.SetDefault(DEFAULT_OPTIONS_DEFAULT_SOURCE, "server1")
	viper.SetDefault(DEFAULT_OPTIONS_THROW_IMMEDIATELY, true)

	// ListOptions defaults
	viper.SetDefault(LIST_DEFAULT_SOURCE, "server2")
	viper.SetDefault(LIST_THROW_IMMEDIATELY, false)
	viper.SetDefault(LIST_MERGE, false)

	// PutOptions defaults
	viper.SetDefault(PUT_DEFAULT_SOURCE, "server1")
	viper.SetDefault(PUT_THROW_IMMEDIATELY, false)
	viper.SetDefault(PUT_CREATE_BUCKET_IF_NOT_EXIST, true)

	// GetObjectOptions defaults
	viper.SetDefault(GET_OBJECT_DEFAULT_SOURCE, "server2")
	viper.SetDefault(GET_OBJECT_THROW_IMMEDIATELY, false)

	// CopyOptions defaults
	viper.SetDefault(COPY_DEFAULT_SOURCE, "server1")
	viper.SetDefault(COPY_THROW_IMMEDIATELY, true)

	// DeleteOptions defaults
	viper.SetDefault(DELETE_DEFAULT_SOURCE, "server1")
	viper.SetDefault(DELETE_THROW_IMMEDIATELY, true)
}
