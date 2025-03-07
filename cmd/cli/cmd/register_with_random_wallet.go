package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(registerWalletCmd)
}

var registerWalletCmd = &cobra.Command{
	Use:   "registerwallet",
	Short: "registers a wallet",
	Args:  cobra.ExactArgs(0),
	Long:  `All software has versions. This is YIP's`,
	Run: func(cmd *cobra.Command, args []string) {

		//c := GetEnvironment(Local)
		//var client = pkg.NewApiClient(fmt.Sprintf("%s/%s", c.yipURL, "api/v1"))
		//
		//fmt.Println("creating wallet...")
		//wallet, err := crypto.GenerateKey()
		//if err != nil {
		//	log.Fatalln(err)
		//	return
		//}
		//key := cryptox.NewKeyFromECDSA(wallet)
		//fmt.Println(key.Address.String())
		//fmt.Println(hexutil.Encode(crypto.FromECDSA(key.PrivateKey)))
		//domain := "http://localhost:3000"
		//_, response, err := client.SIWEChallenge(siwe.ChallengeRequestDTO{
		//	ChainId: "2828",
		//	Address: cryptox.PublicKeyFromPrivateKey(wallet),
		//	Domain:  domain,
		//})
		//
		//if err != nil {
		//	log.Fatalln(err)
		//	return
		//}
		//
		//fmt.Println("login in....")
		//sm, err := cryptox.Sign(response.Challenge, key, cryptox.SignMethodEthereumPrefix, cryptox.SignTypeWeb3JS)
		//
		//if err != nil {
		//	log.Fatalln(err)
		//	return
		//}
		//_, token, err := client.SIWESubmit(siwe.SubmitRequestDTO{
		//	Message:   response.Challenge,
		//	Signature: sm.Signature,
		//	Audience:  "http://localhost:8081",
		//})
		//
		//if err != nil {
		//	log.Fatalln(err)
		//	return
		//}
		//
		//fmt.Println(token.IdToken)

	},
}
