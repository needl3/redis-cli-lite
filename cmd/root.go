package cmd

import (
	"fmt"
	"os"

	"github.com/needl3/redis-cli-lite/client"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "redis-cli-lite",
	Short: "Redis cli lite a simple command line redis client",
	Long:  "This was developed an exercise driven project to learn Golang and it's intricacies.",
	Run: func(cmd *cobra.Command, args []string) {
		host := "localhost"
		port := "6379"
		if len(args) == 2 {
			host = args[0]
			port = args[1]
		}
		redisClient := client.New(host, port)
		redisClient.HandleConnection()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
