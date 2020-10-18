package conf

import (
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"strings"
)

func Parse() *Configuration {
	// Enable configuration through env variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("gs")

	// Set default values
	viper.SetDefault("database.scheme", "postgres")
	viper.SetDefault("database.timezone", "utc")
	viper.SetDefault("database.host", "localhost")
	//viper.SetDefault("database.disable_tls", false)

	// Add config path
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	// Read from configured file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}

	// Unmarshal parsed config to struct
	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}
