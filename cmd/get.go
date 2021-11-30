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
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"

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

func getAssignmentsAll(appConfig appConfig, userId string) (*simpleforce.QueryResult, error) {
	table := "pse__Assignment__c"
	fields := []interface{}{
		"Id",
		"Name",
		"pse__Project__c",
		"pse__Project__r.Name",
		"pse__Project__r.pse__Is_Billable__c",
	}
	query, _, _ := goqu.From(table).Select(fields...).Where(goqu.Ex{
		"pse__Resource__c":                     userId,
		"Open_up_Assignment_for_Time_entry__c": false,
		"pse__Closed_for_Time_Entry__c":        goqu.Op{"eq": true},
	}).ToSql()
	query = strings.ReplaceAll(query, "\"", "")
	query = strings.ReplaceAll(query, "IS", "=")
	result, err := newQuery(appConfig, query)

	if err != nil {
		// handle the error
	}
	return result, nil
}

func getAssignmentsActive(appConfig appConfig, userId string) (*simpleforce.QueryResult, error) {
	table := "pse__Assignment__c"
	fields := []interface{}{
		"Id",
		"Name",
		"pse__Project__c",
		"pse__Project__r.Name",
		"pse__Project__r.pse__Is_Billable__c",
	}
	query, _, _ := goqu.From(table).Select(fields...).Where(goqu.Ex{
		"pse__Resource__c":                     userId,
		"Open_up_Assignment_for_Time_entry__c": false,
		"pse__Closed_for_Time_Entry__c":        false,
		"pse__Exclude_from_Planners__c":        false,
		"pse__End_Date__c":                     goqu.Op{"lt": time.Now().Format("2006-01-02")},
	}).ToSql()

	query = strings.ReplaceAll(query, "\"", "")
	query = strings.ReplaceAll(query, "IS", "=")
	result, err := newQuery(appConfig, query)

	if err != nil {
		// handle the error
	}
	return result, nil
}

func getGlobalProjects(appConfig appConfig) (*simpleforce.QueryResult, error) {
	table := "pse__Proj__c"
	fields := []interface{}{
		"Id",
		"Name",
		"pse__Is_Billable__c",
	}
	query, _, _ := goqu.From(table).Select(fields...).Where(goqu.Ex{
		"pse__Allow_Timecards_Without_Assignment__c": true,
		"pse__Is_Active__c":                          true,
	}).ToSql()

	query = strings.ReplaceAll(query, "\"", "")
	query = strings.ReplaceAll(query, "IS", "=")
	result, err := newQuery(appConfig, query)

	if err != nil {
		// handle the error
	}
	return result, nil
}

func listTimecard(appConfig appConfig, details bool) []string {
	assignments, _ := getAssignmentsAll(appConfig, viper.GetString("userId"))
	table := "pse__Timecard_Header__c"
	fields := []interface{}{
		"Id",
		"Name",
		"pse__Project__c",
		"pse__Assignment__c",
		"pse__Monday_Hours__c",
		"pse__Tuesday_Hours__c",
		"pse__Wednesday_Hours__c",
		"pse__Thursday_Hours__c",
		"pse__Friday_Hours__c",
		"pse__Status__c",
		"OwnerId",
		"PROJECT_ID__c",
		"pse__Approved__c",
		"pse__Start_Date__c",
		"pse__End_Date__c",
		"CreatedById",
		"CreatedDate",
		"IsDeleted",
		"LastModifiedById",
		"LastModifiedDate",
		"LastReferencedDate",
		"LastViewedDate",
		"pse__Audit_Notes__c",
		"pse__Billable__c",
		"pse__Resource__c",
		"pse__Location_Mon__c",
		"pse__Location_Tue__c",
		"pse__Location_Wed__c",
		"pse__Location_Thu__c",
		"pse__Location_Fri__c",
		"pse__Saturday_Hours__c",
		"pse__Saturday_Notes__c",
		"pse__Location_Sat__c",
		"pse__Sunday_Hours__c",
		"pse__Sunday_Notes__c",
		"pse__Location_Sun__c",
		"pse__Timecard_Notes__c",
		"pse__Submitted__c",
		"pse__Monday_Notes__c",
		"pse__Tuesday_Notes__c",
		"pse__Wednesday_Notes__c",
		"pse__Thursday_Notes__c",
		"pse__Friday_Notes__c",
	}

	timeStart := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	timeEnd := time.Now().Format("2006-01-02")

	query, _, _ := goqu.From(table).Select(fields...).Where(goqu.Ex{
		"pse__Start_Date__c": timeStart,
		"pse__End_Date__c":   timeEnd,
		"pse__Resource__c":   viper.GetString("userId"),
	}).ToSql()
	query = strings.ReplaceAll(query, "\"", "")
	query = strings.ReplaceAll(query, "IS", "=")

	queryResult, _ := newQuery(appConfig, query)

	results := []string{}

	for _, record := range queryResult.Records {
		fmt.Println(record)
		for _, assignment := range assignments.Records {
			keys := reflect.ValueOf(assignment).MapKeys()
			for key := range keys {
				fmt.Println(key)
				/* if record.StringField("pse__Assignment__c") == keys[key] {

				} */
			}
		}
	}
	return results
}

func newQuery(appConfig appConfig, query string) (*simpleforce.QueryResult, error) {
	result := &simpleforce.QueryResult{}
	client := simpleforce.NewClient(appConfig.endpoint, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if client == nil {
		return result, fmt.Errorf("")
	}
	err := client.LoginPassword(appConfig.username, appConfig.keychainPassword, appConfig.keychainToken)
	if err != nil {
		return result, fmt.Errorf("")
	}
	result, err = client.Query(query)
	return result, nil
}
