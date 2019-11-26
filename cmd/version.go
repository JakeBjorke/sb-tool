package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//Version is the application version
const Version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print application version",
	Long:  `Print application version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
