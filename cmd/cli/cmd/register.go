package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"yip/pkg"
	"yip/src/api/services/dto"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register admin user",
	Args:  cobra.ExactArgs(3),
	Long:  `All software has versions. This is YIP's`,
	Run: func(cmd *cobra.Command, args []string) {
		location := "https://auth.singularry.xyz"

		password := args[0]
		email := args[1]
		userPassword := args[2]

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

		_, user, err := c.RegisterUser(dto.RegisterRequest{
			Email:    email,
			Password: userPassword,
		})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(user)
	},
}
