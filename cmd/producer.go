package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jakebjorke/sb-tool/sb"
)

// producerCmd represents the producer command
var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Send messages to a service bus",
	Long: `Send messages to a service bus in a configurable 
	fashion`,
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt)

		c := viper.GetString("ConnectionString")
		t := viper.GetString("Producer.Topic")
		intervalStr := viper.GetString("Producer.Upload.Interval")
		dur, err := time.ParseDuration(intervalStr)
		if err != nil {
			fmt.Printf("Unable to parse interval string %s:  %+v\n\r", intervalStr, err)
			return
		}

		ub := viper.GetInt("Producer.Upload.Bytes")

		timeoutStr := viper.GetString("Producer.Upload.Timeout")
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			fmt.Printf("Unable to parse interval string %s:  %+v\n\r", intervalStr, err)
			return
		}

		go sb.Producer(c, t, dur, timeout, ub, done)
		<-done
		fmt.Println()
		fmt.Println("exiting producer")
	},
}

func init() {
	rootCmd.AddCommand(producerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// producerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// producerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
