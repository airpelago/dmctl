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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Shows running containers",
	RunE:  runPs,
}

func runPs(cmd *cobra.Command, args []string) error {
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}
	if len(containers) == 0 {
		bad("No containers running!")
		return nil
	}
	img := viper.GetString("IMAGE")
	if img == "" {
		return errNoImage
	}
	for _, c := range containers {
		if c.Image == imageBase+img {
			good(fmt.Sprintf("Running for %s", time.Since(time.Unix(c.Created, 0)).Truncate(time.Second)))
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(psCmd)
}
