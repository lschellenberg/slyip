package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	infoCmd.AddCommand(infoChain)
}

var infoChain = &cobra.Command{
	Use:   "chain",
	Short: "shows wallet & chain information from slywallet",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("retrieving chain info...")
		cc, err := PrepareApiClient()
		if err != nil {
			log.Fatalln(err)
		}

		_, info, err := cc.apiClient.ChainInfo()
		if err != nil {
			log.Fatalln(err)
		}
		ColoredPrintln("Wallet", info.WalletAddress)
		ColoredPrintln("Wallet Name", info.WalletName)
		ColoredPrintln("Wallet Value", info.WalletValue)
		ColoredPrintln("Chain ID", info.ChainId)
		ColoredPrintln("RPC Location", info.RPCUrl)
	},
}
