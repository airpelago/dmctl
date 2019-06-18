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
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

const (
	droneImage = "docker.io/tobiasfriden/dmc-rpi:latest"
	simImage   = "docker.io/tobiasfriden/dmc-sim:latest"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download latest image versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPullDrone(cmd, args)
	},
}

// pullDrone represents the pull drone command
var pullDrone = &cobra.Command{
	Use:   "drone",
	Short: "Download latest versions of drone image",
	RunE:  runPullDrone,
}

// pullSim represents the pull sim command
var pullSim = &cobra.Command{
	Use:   "sim",
	Short: "Download latest version of sim image",
	RunE:  runPullSim,
}

func runPullDrone(cmd *cobra.Command, args []string) error {
	return pullImage("drone", droneImage)
}

func runPullSim(cmd *cobra.Command, args []string) error {
	return pullImage("sim", simImage)
}

func pullImage(name, imageName string) error {
	fmt.Printf("Pulling %s..\n", name)
	out, err := dockerClient.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	var writer io.Writer
	if Verbose {
		writer = os.Stdout
	} else {
		writer = ioutil.Discard
	}
	_, err = io.Copy(writer, out)
	good("Done!")
	return err
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.AddCommand(pullDrone, pullSim)
}