package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jakebjorke/sb-tool/sb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume messages from a service bus",
	Long: `Consume messages from a service bus in a configurable
	fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt)

		c := viper.GetString("ConnectionString")
		t := viper.GetString("Consumer.Topic")
		s := viper.GetString("Consumer.Subscription")

		timeoutStr := viper.GetString("Consumer.Timeout")
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			fmt.Printf("Unable to parse time out string %s:  %+v\n\r", timeout, err)
			return
		}

		go sb.Consumer(c, t, s, timeout, done)
		<-done
		fmt.Println()
		fmt.Println("exiting producer")
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
