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
	"os"

	"github.com/spf13/cobra"

	"github.com/docker/docker/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var Verbose bool

var dockerClient *client.Client

const dockerFailMessage = `
Could not connect to docker. Docker is necessary in order to run Drone Mission Control onboard software.

If docker is not installed, please install it by following the instructions at: https://docs.docker.com/install/

If docker is running, you might need to run this CLI as sudo. You can also add permissions for the current user to use docker:

  sudo usermod -aG docker USER

Note that these changes require logout to take affect.
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dmctl",
	Short: "CLI to configure and control Drone Mission Control onboard software",
	Long: `CLI to configure and control Drone Mission Control onboard software
  
To perform initial configuration and start the drone software, run:

  dmctl init
  `,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		bad(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Show verbose output")

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Print(dockerFailMessage)
		os.Exit(1)
	}
	dockerClient = cli
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".dmc" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".dmc")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
