package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use: "wygops",
		// Short: "WyGOps is a wireguard server with user based authentification",
		// Long: `With WyGOps it is possible to authenticate users by different
		// 		authentication provider and expose an REST API with the groups
		// 		of the users to handle firewall rules for it`,
		SilenceUsage: true,
		Run:          func(c *cobra.Command, args []string) {},
	}
)

// Execute : Run the cmd parser
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		c.Usage()
		os.Exit(0)
	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
}

func initConfig() {
	viper.AutomaticEnv()

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		fileDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		_configDir := filepath.Join(fileDir, "config")
		viper.AddConfigPath(_configDir)

		readDefaultConfig(_configDir)

		switch env := strings.ToLower(os.Getenv("ENV")); env {
		case "dev":
			viper.SetConfigFile(filepath.Join(_configDir, "dev.yaml"))
			break
		case "prd":
			viper.SetConfigFile(filepath.Join(_configDir, "prd.yaml"))
		case "":
			return
		default:
			log.Fatalf("unknown env variable %s", env)
		}
	}

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("error while reading config %s", err)
	}
}
