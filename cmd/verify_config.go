package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// verify_configCmd represents the verify_config command
var verify_configCmd = &cobra.Command{
	Use:   "verify_config",
	Short: "Read in all configuration and verify it",
	Long:  `Read in all configuration and verify it`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("verify_config called")
	},
}

func init() {
	RootCmd.AddCommand(verify_configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verify_configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verify_configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
