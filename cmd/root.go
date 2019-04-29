package cmd

import (
	"fmt"
	"github.com/larryrun/tor-poller/pkg/config"
	"github.com/larryrun/tor-poller/pkg/poller"
	"github.com/spf13/cobra"
	"os"
)

var configFile string
var rootCmd = &cobra.Command{
	Use:   "tor-poller",
	Short: "tor-poller",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.ReadConfig(configFile)
		if err != nil {
			panic(err)
		}
		poller := poller.NewPoller(conf)
		poller.StartPolling()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "cfg.yaml", "config file")
}