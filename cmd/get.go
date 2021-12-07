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
	"github.com/riweston/gopsa/internal/utils"
	"strings"
	"time"

	"github.com/simpleforce/simpleforce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type timeCardEntry struct {
	projectName      string
	mondayHours      interface{}
	tuesdayHours     interface{}
	wednesdayHours   interface{}
	thursdayHours    interface{}
	fridayHours      interface{}
	saturdayHours    interface{}
	sundayHours      interface{}
	submissionStatus string
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrive timecard related information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		lastWeek, _ := cmd.Flags().GetBool("last-week")
		if lastWeek {
			listTimecard(getConfig(), time.Now().AddDate(0, 0, -7))
		} else {
			listTimecard(getConfig(), time.Now())
		}

		/* 		all, _ := cmd.Flags().GetBool("all")
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
		   		} */
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.Flags().BoolP("last-week", "", false, "Get all assignment")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getAssignmentsAll(appConfig appConfig) (*simpleforce.QueryResult, error) {
	table := "pse__Assignment__c"
	fields := []string{
		"Id",
		"Name",
		"pse__Project__c",
		"pse__Project__r.Name",
		"pse__Project__r.pse__Is_Billable__c",
	}
	filters := []string{
		fmt.Sprintf("pse__Resource__c = '%s'", viper.GetString("userId")),
		"AND",
		"Open_up_Assignment_for_Time_entry__c = false",
		"AND",
		"pse__Closed_for_Time_Entry__c = false",
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ","),
		table,
		strings.Join(filters, " "),
	)
	fmt.Println(query)
	result, err := newQuery(appConfig, query)

	if err != nil {
		// handle the error
	}
	return result, nil
}

func getAssignmentsActive(appConfig appConfig) (*simpleforce.QueryResult, error) {
	table := "pse__Assignment__c"
	fields := []string{
		"Id",
		"Name",
		"pse__Project__c",
		"pse__Project__r.Name",
		"pse__Project__r.pse__Is_Billable__c",
	}
	filters := []string{
		fmt.Sprintf("pse__Resource__c = '%s'", viper.GetString("userId")),
		"AND",
		"Open_up_Assignment_for_Time_entry__c = false",
		"AND",
		"pse__Closed_for_Time_Entry__c = false",
		"AND",
		"pse__Exclude_from_Planners__c = false",
		"AND",
		fmt.Sprintf("pse__End_Date__c = %s", time.Now().Format("2006-01-02")),
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ","),
		table,
		strings.Join(filters, " "),
	)

	result, err := newQuery(appConfig, query)

	if err != nil {
		// handle the error
	}
	return result, nil
}

func getGlobalProjects(appConfig appConfig) (*simpleforce.QueryResult, error) {
	table := "pse__Proj__c"
	fields := []string{
		"Id",
		"Name",
		"pse__Is_Billable__c",
	}
	filters := []string{
		"pse__Allow_Timecards_Without_Assignment__c = true",
		"AND",
		"pse__Is_Active__c",
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ","),
		table,
		strings.Join(filters, " "),
	)

	result, err := newQuery(appConfig, query)

	if err != nil {
		// handle the error
	}
	return result, nil
}

func listTimecard(appConfig appConfig, targetDate time.Time) string {
	table := "pse__Timecard_Header__c"
	fields := []string{
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
	timeStart, timeEnd := utils.DateCalculator(targetDate)
	filters := []string{
		fmt.Sprintf("pse__Start_Date__c = %s", timeStart),
		"AND",
		fmt.Sprintf("pse__End_Date__c = %s", timeEnd),
		"AND",
		fmt.Sprintf("pse__Resource__c = '%s'", viper.GetString("userId")),
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		strings.Join(fields, ","),
		table,
		strings.Join(filters, " "),
	)

	queryResult, _ := newQuery(appConfig, query)
	assignments, _ := getAssignmentsAll(appConfig)
	var result timeCardEntry
	for _, record := range queryResult.Records {
		for _, assignment := range assignments.Records {
			if assignment.StringField("pse__Project__c") == record.StringField("pse__Project__c") {
				result = timeCardEntry{
					projectName:      assignment.StringField("Name"),
					mondayHours:      record.InterfaceField("pse__Monday_Hours__c"),
					tuesdayHours:     record.InterfaceField("pse__Tuesday_Hours__c"),
					wednesdayHours:   record.InterfaceField("pse__Wednesday_Hours__c"),
					thursdayHours:    record.InterfaceField("pse__Thursday_Hours__c"),
					fridayHours:      record.InterfaceField("pse__Friday_Hours__c"),
					saturdayHours:    record.InterfaceField("pse__Saturday_Hours__c"),
					sundayHours:      record.InterfaceField("pse__Sunday_Hours__c"),
					submissionStatus: record.StringField("pse__Status__c"),
				}
			}
			if assignment.StringField("pse__Assignment__c") == record.StringField("pse__Assignment__c") {
				result = timeCardEntry{
					mondayHours:    record.InterfaceField("pse__Monday_Hours__c"),
					tuesdayHours:   record.InterfaceField("pse__Tuesday_Hours__c"),
					wednesdayHours: record.InterfaceField("pse__Wednesday_Hours__c"),
					thursdayHours:  record.InterfaceField("pse__Thursday_Hours__c"),
					fridayHours:    record.InterfaceField("pse__Friday_Hours__c"),
					saturdayHours:  record.InterfaceField("pse__Saturday_Hours__c"),
					sundayHours:    record.InterfaceField("pse__Sunday_Hours__c"),
				}
			}
		}
	}
	
	return fmt.Sprintf("%+v\n", result)
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
