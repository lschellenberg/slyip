package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	infoCmd.AddCommand(infoGetCode)
}

var infoGetCode = &cobra.Command{
	Use:   "codes",
	Short: "shows codes",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("retrieving chain info...")
		cc, err := PrepareApiClient()
		if err != nil {
			log.Fatalln(err)
		}

		_, info, err := cc.apiClient.GetCodes()
		if err != nil {
			log.Fatalln(err)
		}

		for _, v := range info.Data {
			fmt.Printf("%+v\n", v.Code)
		}

	},
}
