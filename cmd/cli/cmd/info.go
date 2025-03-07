package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"yip/pkg"
	"yip/src/api/services/dto"
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
		env := GetEnvironment(Cloud)
		location := env.yipURL
		password := args[0]
		fmt.Println(password)

		c := pkg.NewApiClient(fmt.Sprintf("%s/%s", location, "api/v1"))
		_, r, err := c.SIWEChallenge(dto.ChallengeRequestDTO{
			ChainId: "11155111",
			Address: "0x01200120310230123123",
			Domain:  "https://domain.com",
		})

		if err != nil {
			log.Fatalln(err)
		}

		b, err := json.Marshal(r)

		fmt.Println(r)
		fmt.Println(len(b))
		//_, token, err := c.SignIn(services.SignInRequest{
		//	Email:    "leonard.schellenberg@gmail.com",
		//	Password: password,
		//	Audiences: []string{
		//		location,
		//	},
		//})
		//
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//
		//c.SetToken(token.IdToken)
		//
		//_, users, err := c.GetUsers()
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//for _, u := range users.Users {
		//	fmt.Println(u.Email, u.Role, u.ID, u.IsAdmin())
		//}
	},
}
