/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/simpleforce/simpleforce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrive timecard related information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		if all {
			data, _ := getAssignmentsAll(getConfig(), viper.GetString("userId"))
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			for _, v := range data {
				t.AppendRows([]table.Row{{v.Name}})
			}
			t.Render()
		} else {
			data, _ := getAssignmentsActive(getConfig(), viper.GetString("userId"))
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)

			for _, v := range data {
				t.AppendRows([]table.Row{{v.Name}})
			}
			t.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.Flags().BoolP("all", "", false, "Get all assignment")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type projectAssignment struct {
	Id               string
	Name             string
	Project          string
	Project_Name     string
	Project_Billable string
}

type projectGlobalAssignment struct {
	Project_Id   string
	Project_Name string
	Billable     string
}

func getAssignmentsAll(appConfig appConfig, userId string) ([]projectAssignment, error) {
	var list []projectAssignment
	fields := "Id, Name, pse__Project__c, pse__Project__r.Name, pse__Project__r.pse__Is_Billable__c"
	filters := fmt.Sprintf("pse__Resource__c = '%s' AND Open_up_Assignment_for_Time_entry__c = false AND pse__Closed_for_Time_Entry__c = false", userId)
	query := "SELECT " + fields + " FROM pse__Assignment__c " + " WHERE " + filters

	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		// handle the error

		return list, fmt.Errorf("")
	}

	err := client.LoginPassword(appConfig.username, appConfig.keychainPassword, appConfig.keychainToken)
	if err != nil {
		// handle the error

		return list, fmt.Errorf("")
	}
	result, err := client.Query(query)
	for _, record := range result.Records {
		list = append(list, projectAssignment{
			Id:               record.StringField("Id"),
			Name:             strings.TrimRight(record.StringField("Name"), "\r\n"),
			Project:          record.StringField("pse__Project__c"),
			Project_Name:     record.StringField("pse__Project__r.Name"),
			Project_Billable: record.StringField("pse__Project__r.pse__Is_Billable__c"),
		})
	}
	if err != nil {
		// handle the error
	}
	return list, nil
}

func getAssignmentsActive(appConfig appConfig, userId string) ([]projectAssignment, error) {
	var list []projectAssignment
	currentDate := time.Now().Format("2006-01-02")
	fields := "Id, Name, pse__Project__c, pse__Project__r.Name, pse__Project__r.pse__Is_Billable__c"
	filters := fmt.Sprintf("pse__Resource__c = '%s' AND Open_up_Assignment_for_Time_entry__c = false AND pse__Closed_for_Time_Entry__c = false AND pse__Exclude_from_Planners__c = false AND pse__End_Date__c > %s", userId, currentDate)
	query := "SELECT " + fields + " FROM pse__Assignment__c " + " WHERE " + filters

	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		// handle the error

		return list, fmt.Errorf("")
	}

	err := client.LoginPassword(appConfig.username, appConfig.keychainPassword, appConfig.keychainToken)
	if err != nil {
		// handle the error

		return list, fmt.Errorf("")
	}
	result, err := client.Query(query)
	for _, record := range result.Records {
		list = append(list, projectAssignment{
			Id:               record.StringField("Id"),
			Name:             strings.TrimRight(record.StringField("Name"), "\r\n"),
			Project:          record.StringField("pse__Project__c"),
			Project_Name:     record.StringField("pse__Project__r.Name"),
			Project_Billable: record.StringField("pse__Project__r.pse__Is_Billable__c"),
		})
	}
	if err != nil {
		// handle the error
	}
	return list, nil
}

func getGlobalProjects(appConfig appConfig) ([]projectGlobalAssignment, error) {
	var list []projectGlobalAssignment
	fields := "Id, Name, pse__Is_Billable__c"
	filters := "pse__Allow_Timecards_Without_Assignment__c = true and pse__Is_Active__c = true"
	query := "SELECT " + fields + " FROM pse__Proj__c " + " WHERE " + filters

	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		// handle the error

		return list, fmt.Errorf("")
	}

	err := client.LoginPassword(appConfig.username, appConfig.keychainPassword, appConfig.keychainToken)
	if err != nil {
		// handle the error

		return list, fmt.Errorf("")
	}
	result, err := client.Query(query)
	for _, record := range result.Records {
		list = append(list, projectGlobalAssignment{
			Project_Id:   record.StringField("Id"),
			Project_Name: strings.TrimRight(record.StringField("Name"), "\r\n"),
			Billable:     record.StringField("pse__Project__c"),
		})
	}
	if err != nil {
		// handle the error
	}
	return list, nil
}
