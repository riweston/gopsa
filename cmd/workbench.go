package main

import (
	"fmt"
	"strings"

	"github.com/simpleforce/simpleforce"
)

const (
)

type AppConfig struct {
	endpoint string
	username string
	password string
	token    string
}

func setConfig(endpoint string, username string, password string, token string) (appConfig AppConfig) {
	appConfig = AppConfig{
		endpoint: endpoint,
		username: username,
		password: password,
		token:    token,
	}
	return
}

func getUserId(appConfig AppConfig) (userId string) {
	/* 	endpoint := os.Getenv("endpoint")
	   	username := os.Getenv("username")
	   	password := os.Getenv("password")
	   	token := os.Getenv("token") */

	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		// handle the error

		return
	}

	err := client.LoginPassword(appConfig.username, appConfig.password, appConfig.token)
	if err != nil {
		// handle the error

		return
	}

	query := "select Id, Name, Email from Contact where pse__Is_Resource__c = true and Email LIKE '" + username + "'"
	result, err := client.Query(query)
	if err != nil {
		// handle the error
	}
	for _, record := range result.Records {
		userId = record.StringField("Id")
	}
	return
}

func getAssignmentsAll(appConfig AppConfig, userId string) {
	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		// handle the error

		return
	}

	err := client.LoginPassword(appConfig.username, appConfig.password, appConfig.token)
	if err != nil {
		// handle the error

		return
	}
	filter := []string{
		"select Id",
		"Name",
		"pse__Project__c",
		"pse__Project__r.Name",
		"pse__Project__r.pse__Is_Billable__c from pse__Assignment__c where pse__Resource__c = '" + userId + "' and Open_up_Assignment_for_Time_entry__c = false and pse__Closed_for_Time_Entry__c = false",
	}
	query := strings.Join(filter, ", ")
	result, err := client.Query(query)
	if err != nil {
		// handle the error
	}
	for _, record := range result.Records {
		userId = record.StringField("Id")
	}
	return
}

func main() {
	config := setConfig(endpoint, username, password, token)
	storeConfig(config)
	userId := getUserId(config)
	fmt.Print(userId)
}
