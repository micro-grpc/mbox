// Copyright © 2018 Oleg Dolya <oleg.dolya@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"strings"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate JWT token",
	Long:  `Generate JWT token.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		secret := viper.GetString("secret")
		if len(secret) == 0 {
			log.Fatalln("IS NOT SET SECRET")
		}
		// TokenAuth = jwtauth.New("HS256", []byte(viper.GetString("secret")), nil)
		ClientPrompts = []*survey.Question{
			{
				Name:      "name",
				Prompt:    &survey.Input{Message: "User"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
    tokenAuth := jwtauth.New("HS256", []byte(viper.GetString("secret")), nil)
		items := jwtauth.Claims{}
		if len(args) > 0 {
			// For debugging/example purposes, we generate and print
			// a sample jwt token with claims `user:test` here:

			cf, err := cmd.Flags().GetStringSlice("claims")
			if err != nil {
				log.Fatalln(err.Error())
			}

			items["user_id"] = args[0]
			items["user"] = args[0]
			user, err := cmd.Flags().GetString("name")
			if err == nil {
				if len(user) > 0 {
					items["user"] = user
				}
			}
			fmt.Printf("User ID: %s\n", items["user_id"])
			fmt.Printf("User: %s\n", items["user"])

			if len(cf) > 0 {
				fmt.Println("Other claims:")
				for _, v := range cf {
					item := strings.Split(v, "=")
					if len(item) == 2 {
						fmt.Printf("%v: %v\n", item[0], item[1])
						items[item[0]] = item[1]
					}
				}
			}

			_, tokenString, _ := tokenAuth.Encode(items)
			fmt.Printf("\nAuthorization: BEARER %s\n\n", tokenString)
			fmt.Printf("export TOKEN=%s\n\n", tokenString)

			// fmt.Printf("Secret: %s\n", viper.GetString("secret"))
		} else {
			answers := struct {
				Name          string // survey will match the question and field names
				FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
				Age           int    // if the types don't match exactly, survey will try to convert for you
			}{}

			err := survey.Ask(ClientPrompts, &answers)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			items["user_id"] = strings.ToLower(answers.Name)
			items["user"] = strings.ToLower(answers.Name)
			_, tokenString, _ := tokenAuth.Encode(items)
			fmt.Printf("\nAuthorization: BEARER %s\n\n", tokenString)
			fmt.Printf("export TOKEN=%s\n\n", tokenString)
		}
	},
}

func init() {
	configCmd.AddCommand(tokenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	tokenCmd.Flags().StringP("name", "n", "", "user name if not equal user_id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tokenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tokenCmd.Flags().StringSliceP("claims", "c", []string{}, "Claims")
}
