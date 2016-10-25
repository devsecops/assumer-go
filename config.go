package assumer

import "github.com/spf13/viper"

// Config initializes the assumer configuration
func Config() error {
	viper.SetConfigName("assumer")        // name of config file (without extension)
	viper.AddConfigPath("$HOME/.assumer") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		return err
	}
	return nil
}
