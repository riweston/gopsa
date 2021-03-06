package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/simpleforce/simpleforce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "User configuration for gopsa",
	Long:  `Use to configure the gopsa CLI application with your implementation of PSA`,
	Run: func(cmd *cobra.Command, args []string) {
		setConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func setConfig() {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("Nothing entered")
		}
		return nil
	}

	endpoint := promptui.Prompt{
		Label:    "Enter Salesforce URL",
		Validate: validate,
	}

	username := promptui.Prompt{
		Label:    "Enter Salesforce username",
		Validate: validate,
	}

	password := promptui.Prompt{
		Label:    "Enter Salesforce password",
		Validate: validate,
		Mask:     '*',
	}

	token := promptui.Prompt{
		Label:    "Enter Salesforce API token",
		Validate: validate,
		Mask:     '*',
	}

	resultEndpoint, err := endpoint.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	resultUsername, err := username.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	resultPassword, err := password.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	resultToken, err := token.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	setConfigEndpoint(resultEndpoint)
	setConfigUsername(resultUsername)
	setConfigPassword(resultEndpoint, resultUsername, resultPassword)
	setConfigToken(resultEndpoint, resultUsername, resultToken)
	appConfig := appConfig{
		endpoint:         resultEndpoint,
		username:         resultUsername,
		keychainPassword: resultPassword,
		keychainToken:    resultToken,
	}
	setConfigUserId(appConfig)
}

func setConfigEndpoint(endpoint string) {
	viper.Set("endpoint", endpoint)
	viper.WriteConfigAs(viper.ConfigFileUsed())
}

func setConfigUsername(username string) {
	viper.Set("username", username)
	viper.WriteConfigAs(viper.ConfigFileUsed())
}

func setConfigPassword(endpoint string, username string, password string) {
	service := "gopsa." + endpoint + ".password"
	err := keyring.Set(service, username, password)
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("keychainPassword", service)
	viper.WriteConfigAs(viper.ConfigFileUsed())
}

func setConfigToken(endpoint string, username string, token string) {
	service := "gopsa." + endpoint + ".token"
	err := keyring.Set(service, username, token)
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("keychainToken", service)
	viper.WriteConfigAs(viper.ConfigFileUsed())
}

func setConfigUserId(appConfig appConfig) {
	viper.Set("userId", getUserId(appConfig))
	viper.WriteConfigAs(viper.ConfigFileUsed())
}

type appConfig struct {
	endpoint         string
	username         string
	keychainPassword string
	keychainToken    string
}

func getConfig() appConfig {
	password, _ := keyring.Get(viper.GetString("keychainPassword"), viper.GetString("username"))
	token, _ := keyring.Get(viper.GetString("keychainToken"), viper.GetString("username"))
	return appConfig{
		endpoint:         viper.GetString("endpoint"),
		username:         viper.GetString("username"),
		keychainPassword: password,
		keychainToken:    token,
	}
}

func getUserId(appConfig appConfig) (userId string) {
	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		// handle the error

		return
	}

	err := client.LoginPassword(appConfig.username, appConfig.keychainPassword, appConfig.keychainToken)
	if err != nil {
		// handle the error

		return
	}
	fields := "Id, Name, Email"
	filters := "pse__Is_Resource__c = true AND Email LIKE '" + appConfig.username + "'"
	query := "SELECT " + fields + " FROM Contact " + " WHERE " + filters
	result, err := client.Query(query)
	if err != nil {
		// handle the error
	}
	for _, record := range result.Records {
		userId = record.StringField("Id")
	}
	return
}
