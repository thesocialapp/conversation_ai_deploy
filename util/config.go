package util

import "github.com/spf13/viper"

type Config struct {
	RedisAddr            string `mapstructure:"REDIS_ADDRESS"`
	HttpServerAddress    string `mapstructure:"HTTP_SERVER_ADDRESS"`
	RedisPrefix          string `mapstructure:"REDIS_PREFIX"`
	SocketIOPingTimeout  int    `mapstructure:"SOCKET_IO_PING_TIMEOUT"`
	SocketIOPingInterval int    `mapstructure:"SOCKET_IO_PING_INTERVAL"`
	OpenAPIKey           string `mapstructure:"OPEN_API_KEY"`
}

// Load the app.env file and unmarshal it into the Config struct
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
