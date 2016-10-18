package assumer

import (
	"fmt"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration
func InitConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.AddConfigPath("$HOME/.assumer") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
