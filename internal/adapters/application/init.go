package application

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/AcroManiac/micropic/internal/adapters/logger"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init(defaultConfigPath string) {
	// using standard library "flag" package
	flag.String("config", defaultConfigPath, "path to configuration flag")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	// Check config file existence and read configuration from it
	configPath := viper.GetString("config") // retrieve value from viper
	if _, err := os.Stat(configPath); err == nil {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Couldn't read configuration file: %s", err.Error())
		}
	} else {
		log.Printf("File %v doesn't exist. Trying to read from environment variables", configPath)
		viper.SetConfigType("env")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	// Setting log parameters
	logger.Init(viper.GetString("log.level"), viper.GetString("log.file"))
}
