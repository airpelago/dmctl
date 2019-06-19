/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
)

var (
	Follow bool
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs from running containers",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runLogs,
}

func runLogs(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return containerLogs("/" + args[0])
	}
	return containerLogs("/drone")

}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.PersistentFlags().BoolVarP(&Follow, "follow", "f", false, "Attach and continously output logs")
}
