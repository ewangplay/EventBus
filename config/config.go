package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config parse eventbus service configure
type Config struct {
	viper *viper.Viper
	cfg   map[string]interface{}
}

// New instance of Config
func New(filename string) (config *Config, err error) {
	config = &Config{}

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	lviper := viper.New()

	fName := filepath.Base(filename)
	extName := filepath.Ext(filename)
	cmdRoot := fName[:len(fName)-len(extName)]
	lviper.SetEnvPrefix(strings.ToUpper(cmdRoot))
	lviper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	lviper.SetEnvKeyReplacer(replacer)

	// Now set the configuration file.
	lviper.SetConfigName(cmdRoot)                // Name of config file (without extension)
	lviper.AddConfigPath(filepath.Dir(filename)) // Path to look for the config file in

	err = lviper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, err
	}

	config.viper = lviper

	return config, nil
}

// GetConfig returns K-V structure configure data
func (c *Config) GetConfig() map[string]interface{} {
	if c.cfg != nil {
		return c.cfg
	}

	if c.viper != nil {
		c.cfg = make(map[string]interface{})

		c.convertSettings("", c.viper.AllSettings(), c.cfg)

		return c.cfg
	}
	return nil
}

func (c *Config) convertSettings(parent string, src map[string]interface{}, dest map[string]interface{}) {
	var key string
	for k, v := range src {
		key = k
		if parent != "" {
			key = parent + "_" + k
		}

		mv, ok := v.(map[string]interface{})
		if ok {
			c.convertSettings(key, mv, dest)
		} else {
			dest[key] = v
		}
	}
}
