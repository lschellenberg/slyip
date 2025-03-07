package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"yip/pkg"
	"yip/src/api/services/dto"
)

func init() {
	rootCmd.AddCommand(setEmailCmd)
}

var setEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Set Email Of User",
	Args:  cobra.ExactArgs(3),
	Long:  `Sets the email of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		env := GetEnvironment(Local)
		location := env.yipURL

		password := args[0]
		email := args[1]
		userId := args[2]

		c := pkg.NewApiClient(fmt.Sprintf("%s/%s", location, "api/v1"))

		_, token, err := c.SignIn(dto.SignInRequest{
			Email:    "leonard.schellenberg@gmail.com",
			Password: password,
			Audiences: []string{
				location,
			},
		})

		if err != nil {
			log.Fatalln(err)
		}

		c.SetToken(token.IdToken)

		fmt.Println("Setting email")
		status, user, err := c.SetEmail(dto.SetEmailRequest{
			Email:  email,
			UserID: userId,
		})
		fmt.Println(status)

		if err != nil {
			log.Fatalln(err)
			return
		}
		fmt.Println(user)
	},
}
