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
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
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

func containerLogs(name string) error {
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	var id string
	for _, c := range containers {
		for _, n := range c.Names {
			if n == name {
				id = c.ID
				break
			}
		}
	}
	if id == "" {
		bad(fmt.Sprintf("Container %s not found", strings.TrimPrefix(name, "/")))
		return nil
	}
	out, err := dockerClient.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     Follow,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	return err
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.PersistentFlags().BoolVarP(&Follow, "follow", "f", false, "Attach and continously output logs")
}
