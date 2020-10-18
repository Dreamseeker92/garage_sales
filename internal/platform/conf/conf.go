package conf

import (
	"fmt"
	"time"
)

type Configuration struct {
	Server   *ServerConfiguration   `yaml:"server"`
	Database *DatabaseConfiguration `yaml:"database"`
}

type ServerConfiguration struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

type DatabaseConfiguration struct {
	Scheme     string `yaml:"scheme"`
	Timezone   string `yaml:"timezone"`
	Host       string `yaml:"host"`
	Name       string `yam:"name"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	DisableTLS bool   `mapstructure:"disable_tls"`
}

func (config *Configuration) Address() string {
	return fmt.Sprintf("%s:%v", config.Server.Host, config.Server.Port)
}

func (config *Configuration) ReadTimeout() time.Duration {
	return time.Duration(config.Server.ReadTimeout) * time.Second
}

func (config *Configuration) WriteTimeout() time.Duration {
	return time.Duration(config.Server.ReadTimeout) * time.Second
}

func (config *Configuration) ShutdownTimeout() time.Duration {
	return time.Duration(config.Server.ShutdownTimeout) * time.Second
}
