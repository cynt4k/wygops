package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		c.Usage()
		os.Exit(0)
	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	// rootCmd.PostRun(func(c *cobra.Command, args []string) {
	// 	os.Exit(0)
	// })
}
