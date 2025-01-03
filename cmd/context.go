package cmd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

const (
	configFileName = "unpost"
)

var Token string

var contextCommand = &cobra.Command{
	Use:   "context",
	Short: "context commands",
}

func init() {
	contextCommand.AddCommand(setContextCommand())
	contextCommand.AddCommand(currentContextCommand())
	contextCommand.AddCommand(resetContextCommand())
}

type Context struct {
	Token string `json:"token"`
}

// saves the context info to the config file in ~/.config/authbase
func setContextCommand() *cobra.Command {
	var token string
	var organization string
	command := &cobra.Command{
		Use:   "set",
		Short: "set context",
		Run: func(cmd *cobra.Command, args []string) {
			if token == "" || organization == "" {
				logrus.Infof(`missing required flags: --token, --organization`)
				return
			}

			// save the context info to the config file
			viper.SetConfigName("unpost")
			viper.AddConfigPath("./.tmp")
			viper.SetConfigType("yml")
			viper.Set("context", Context{
				Token: token,
			})

			if err := viper.WriteConfig(); err != nil {
				fmt.Println("error writing config file: ", err)
			} else {
				fmt.Println("context saved")
			}
		},
	}

	command.Flags().StringVarP(&token, "token", "t", "", "token")

	return command
}

func currentContextCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "current",
		Short: "current context",
	}

	return command
}

func resetContextCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "reset",
		Short: "reset context",
	}

	return command
}

func verifyContext() {
	if Token == "" {
		logrus.Error("missing required flags: --token")
		return
	}
}

func writeContext(context Context) {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath("./.tmp")
	viper.SetConfigType("yml")
	viper.Set("context", context)

	if err := viper.WriteConfig(); err != nil {
		fmt.Println("error writing config file: ", err)
	}
}

func bindContextFlags(command *cobra.Command) {
	command.Flags().StringVarP(&Token, "token", "t", "", "token")
}

func tokenContext() context.Context {
	cfg := readContext()
	Token = cfg.Token

	md := metadata.New(map[string]string{"Authorization": "Bearer " + Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}

func readContext() Context {
	var context Context
	viper.SetConfigName(configFileName)
	viper.AddConfigPath("./.tmp")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading config file: ", err)
	}

	if err := viper.UnmarshalKey("context", &context); err != nil {
		fmt.Println("error unmarshalling config file: ", err)
	}

	return context
}
