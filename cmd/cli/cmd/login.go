package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"yip/pkg"
	"yip/src/api/services/dto"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login admin user",
	Args:  cobra.ExactArgs(3),
	Long:  `All software has versions. This is YIP's`,
	Run: func(cmd *cobra.Command, args []string) {
		location := "https://auth.singularry.xyz"

		email := args[0]
		userPassword := args[1]
		audience := args[2]

		c := pkg.NewApiClient(fmt.Sprintf("%s/%s", location, "api/v1"))

		_, token, err := c.SignIn(dto.SignInRequest{
			Email:    email,
			Password: userPassword,
			Audiences: []string{
				audience,
			},
		})

		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(token.IdToken)
		fmt.Println(token.RefreshToken)
	},
}
