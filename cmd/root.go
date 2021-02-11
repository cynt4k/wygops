package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	c = config.GetConfig()
)

// Execute : Run the cmd parser
func Execute() error {
	return rootCmd.Execute()
}

func init() { // nolint:gochecknoinits
	cobra.OnInitialize(initConfig)
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		err := c.Usage()

		if err != nil {
			log.Fatal(err)
		}
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

		err = readDefaultConfig(_configDir, c)

		if err != nil {
			log.Fatal(err)
		}

		switch env := strings.ToLower(os.Getenv("ENV")); env {
		case "dev":
			viper.SetConfigFile(filepath.Join(_configDir, "dev.yaml"))
		case "prd":
			viper.SetConfigFile(filepath.Join(_configDir, "prd.yaml"))
		case "":
			if err := viper.Unmarshal(&c); err != nil {
				log.Fatalf("error while unmarshal config %s", err)
			}
			return
		default:
			log.Fatalf("unknown env variable %s", env)
		}
	}

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("error while reading config %s", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("error while unmarshal config %s", err)
	}
	if os.Getenv("ENV") == "dev" {
		viper.Set("DevMode", true)
		c.DevMode = true
	}
}

func getLogger() (logger *zap.Logger) {
	if c.DevMode {
		return getCLILogger()
	}
	return getCLILogger()
}

func getCLILogger() (logger *zap.Logger) {
	var level zap.AtomicLevel

	if c.DevMode {
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		level = zap.NewAtomicLevel()
	}
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cfg := zap.Config{
		Level:       level,
		Development: c.DevMode,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ = cfg.Build()
	return logger
}
