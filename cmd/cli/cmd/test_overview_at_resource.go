package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"yip/pkg"
	"yip/src/api/services/dto"
)

func init() {
	rootCmd.AddCommand(testOverviewAtResourceCmd)
}

var testOverviewAtResourceCmd = &cobra.Command{
	Use:   "testoverview",
	Short: "delivers overview of test",
	Args:  cobra.ExactArgs(2),
	Long:  `All software has versions. This is YIP's`,
	Run: func(cmd *cobra.Command, args []string) {

		email := args[0]
		userPassword := args[1]

		c := GetEnvironment(Cloud)

		var client = pkg.NewApiClient(fmt.Sprintf("%s/%s", c.yipURL, "api/v1"))

		_, token, err := client.SignIn(dto.SignInRequest{
			Email:    email,
			Password: userPassword,
			Audiences: []string{
				c.oracleURL,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		resourceClient := pkg.NewHttpClient(c.oracleURL)
		var users *GetUsersResponse = &GetUsersResponse{}
		status, err := resourceClient.Get(users, token.IdToken, "admin/users")
		if err != nil {
			log.Fatal(err)
		}

		if status >= 300 {
			log.Fatal(err)
		}

		filteredUsers := make([]User, 0)

		for _, u := range users.Data {
			if isInWhiteList(u.Name) {
				filteredUsers = append(filteredUsers, u)
			}
		}

		for _, u := range filteredUsers {
			fmt.Println("================================================================================================================")
			fmt.Println(u.Name, u.ID)
			showChoices(u.ID, &resourceClient, token.IdToken)
		}

	},
}

func showChoices(id string, client *pkg.HttpClient, token string) {

	var choices *GetChoices = &GetChoices{}
	status, err := client.Get(choices, token, fmt.Sprintf("admin/users/%s/choices", id))
	if err != nil {
		log.Fatal(err)
	}

	if status >= 300 {
		log.Fatal(err)
	}

	likes := 0
	dislikes := 0
	unvoted := 0

	for _, v := range choices.Choices {
		switch v.LikeType {
		case 0:
			unvoted++
		case 1:
			likes++
		case 2:
			dislikes++
		}
	}

	fmt.Println("Likes:", likes, ", Dislikes:", dislikes, ", Not Voted Yet:", unvoted)

}

func isInWhiteList(name string) bool {
	for _, w := range whitelist {
		if w == name {
			return true
		}
	}
	return false
}

var whitelist = []string{
	"Haris",
	"Lenny",
	"Emily",
	"Sam",
	"Jen Patryn",
	"George",
	"Thomas Bailey",
}

/*
Haris 0x437ec6809fEe4CB7c31112826a209D3d9bb86a9B
Haris 0x101e1f15B59eBb1767DD0cFf0EbE3BbBd0b02CB1
Haris 0x703dEab53CC765b839BDd382902901C42522d638
Lenny 0x2B978c99Fc539818E3fB4Fe0F2c216f78D78eF16
Haris 0xd37696E645651C08BF000AA2e9d253Fe7c06185C
Lenny 0x8aE1C731c92b21a50f2AB76c66cf18a4182220df
Haris 0x300E4477120BC0375644Df2342e7954911163Acc
Lenny 0x7E1ce2A23d1813a10B1291283399Cd7bb48aFd06
Haris 0x70eF5D4402D1453752439b1bE6aE6636a009dBFd
Lenny 0x76A3ecC5B24c131a1Dc6C7ebf0f301eb629fF5ee
Emily 0x004731eaB8E8BF8C9aEb67BF2DDe6Ebeb76C1B6C
Sam 0xCbeF5CF800C70AC99f41ddcA72e967b5B9500051
Jen Patryn 0x9560d7b0a67774F434A1B540fF97eF88428897Ef
George 0x5CF6714903167E65e372e3F5041F45616693E2F8
Thomas Bailey 0x9549eD7dA0ec92E6c7D0f9DBfaF8354D2BDD5C47

*/
