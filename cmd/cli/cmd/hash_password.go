package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"yip/src/cryptox"
)

func init() {
	rootCmd.AddCommand(hashPasswordCmd)
}

var hashPasswordCmd = &cobra.Command{
	Use:   "hash",
	Short: "hashes a password",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		password := args[0]
		hashed, err := cryptox.HashPassword(password)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(hashed)
	},
}
