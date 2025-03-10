package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the info of YIP",
	Args:  cobra.ExactArgs(1),
	Long:  `All software has versions. This is YIP's`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
