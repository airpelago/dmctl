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
	"github.com/spf13/viper"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop dmc containers",
	RunE:  runStopDrone,
}

// stopCmd represents the stop command
var stopDroneCmd = &cobra.Command{
	Use:   "drone",
	Short: "Stop drone container",
	RunE:  runStopDrone,
}

func runStopDrone(cmd *cobra.Command, args []string) error {
	img := viper.GetString("IMAGE")
	if img == "" {
		return errNoImage
	}
	return runStop("drone", img)
}

func runStop(name, imageName string) error {
	running, err := containerRunning(imageName)
	if err != nil {
		return err
	}
	if !running {
		bad(name + " not running")
		return nil
	}
	return stopContainer(name, imageName)
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.AddCommand(stopDroneCmd)
}
