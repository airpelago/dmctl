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
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Recreate  bool
	NoRestart bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start dmc containers",
	RunE:  runStartDrone,
}

// startDroneCmd represents the start drone command
var startDroneCmd = &cobra.Command{
	Use:   "drone",
	Short: "Start drone container",
	RunE:  runStartDrone,
}

func runStartDrone(cmd *cobra.Command, args []string) error {
	img := viper.GetString("IMAGE")
	if img == "" {
		return errNoImage
	}
	if img == "dmc-sim" {
		return runSimulatedDrone(img)
	} else {
		return runOnboardDrone(img)
	}
}

func runOnboardDrone(imageName string) error {
	droneEnv := envList("ID", "PASSWORD", "FCU_URL", "GCS_URL", "DMC_URI", "DMC_SESSION_URI", "DMC_ANIP_URI", "MOCK_IMSI", "MOCK_POSITION")
	config := &container.Config{
		Env: droneEnv,
		Cmd: []string{},
		Tty: true,
	}
	policy := container.RestartPolicy{}
	if !NoRestart {
		policy.Name = "unless-stopped"
	}
	hostConfig := &container.HostConfig{
		Privileged:    true,
		NetworkMode:   "host",
		RestartPolicy: policy,
	}
	return startContainer("drone", imageName, config, hostConfig)
}

func runSimulatedDrone(imageName string) error {
	location := viper.GetString("MOCK_POSITION")
	if location == "" {
		return errors.New("location must be set for simulation")
	}
	if err := writeConfig(); err != nil {
		return errors.New("failed writing location to config")
	}
	simType := viper.GetString("SIM_TYPE")
	if simType == "" {
		return errors.New("simulation type not set, run dmctl config obc")
	}
	droneEnv := envList("ID", "PASSWORD", "DMC_URI", "DMC_SESSION_URI", "DMC_ANIP_URI", "MOCK_IMSI", "MOCK_POSITION")
	config := &container.Config{
		Env: droneEnv,
		Cmd: []string{
			fmt.Sprintf("--location %s,0", location),
			fmt.Sprintf("--%s", simType),
		},
		Tty: true,
	}
	hostConfig := &container.HostConfig{}
	return startContainer("drone", imageName, config, hostConfig)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(startDroneCmd)

	startCmd.PersistentFlags().BoolVarP(&Recreate, "recreate", "r", false, "Recreate if already running")
	startCmd.Flags().StringP("location", "l", "", "Simulation location (LAT,LNG,ALT)")
	if err := viper.BindPFlag("MOCK_POSITION", startCmd.Flags().Lookup("location")); err != nil {
		panic(err)
	}

	startDroneCmd.PersistentFlags().BoolVarP(&NoRestart, "no_restart", "n", false, "Do not enable automatic restart")

}
