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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	httpClient = http.DefaultClient
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure dmc settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigureDrone(cmd, args)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to authorize with dmc",
	RunE: func(dmc *cobra.Command, args []string) error {
		token, err := login()
		if err != nil {
			return err
		}
		viper.Set("TOKEN", token)
		return writeConfig()
	},
}

var droneCmd = &cobra.Command{
	Use:   "drone",
	Short: "Configures which registered drone to connect",
	RunE:  runConfigureDrone,
	PostRun: func(cmd *cobra.Command, args []string) {
		good("Drone config saved!")
	},
}

var obcCmd = &cobra.Command{
	Use:   "obc",
	Short: "Configure onboard software version",
	RunE:  rungConfigureOBC,
}

var anipCmd = &cobra.Command{
	Use:   "anip",
	Short: "Configures ANIP settings",
	RunE:  runConfigureANIP,
	PostRun: func(cmd *cobra.Command, args []string) {
		good("ANIP config saved!")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all configuration",
	RunE:  runConfigList,
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears all configuration",
	RunE:  runConfigClear,
}

func runConfigureDrone(cmd *cobra.Command, args []string) error {
	t, err := token()
	if err != nil {
		return err
	}
	drones, err := getDrones(t)
	if err != nil {
		return err
	}
	id, pass, url, err := droneConfig(drones)
	if err != nil {
		return err
	}
	viper.Set("ID", id)
	viper.Set("PASSWORD", pass)
	viper.Set("FCU_URL", url)
	return writeConfig()
}

var images = []string{
	"dmc-rpi",
	"dmc-x86",
	"dmc-sim",
}

func rungConfigureOBC(cmd *cobra.Command, args []string) error {
	obcPrompt := &promptui.Select{
		Label: "Select OBC type",
		Items: []string{
			"Raspberry Pi",
			"Linux 64-bit",
			"Simulated",
		},
	}
	idx, _, err := obcPrompt.Run()
	if err != nil {
		return err
	}
	viper.Set("IMAGE", images[idx])
	return writeConfig()
}

func runConfigureANIP(cmd *cobra.Command, args []string) error {
	uriPrompt := &promptui.Prompt{
		Label: "ANIP uri",
	}
	uri, err := uriPrompt.Run()
	if err != nil {
		return err
	}
	viper.Set("DMC_ANIP_URI", uri)

	imsiPrompt := &promptui.Prompt{
		Label: "Mock IMSI",
	}
	imsi, err := imsiPrompt.Run()
	if err != nil {
		return err
	}
	viper.Set("MOCK_IMSI", imsi)

	posPrompt := &promptui.Prompt{
		Label: "Mock position (LAT,LNG,ALT)",
	}
	pos, err := posPrompt.Run()
	if err != nil {
		return err
	}
	viper.Set("MOCK_POSITION", pos)

	return writeConfig()
}

func runConfigList(cmd *cobra.Command, args []string) error {
	dir, err := homedir.Dir()
	if err != nil {
		return err
	}
	path := dir + "/.dmc.yaml"
	if raw, err := ioutil.ReadFile(path); err != nil {
		if os.IsNotExist(err) {
			bad("No config file created")
			return nil
		}
		return err
	} else {
		var config map[string]interface{}
		if err := yaml.Unmarshal(raw, &config); err != nil {
			return err
		}
		if _, ok := config["token"]; ok {
			delete(config, "token")
		}
		filtered, err := yaml.Marshal(&config)
		if err != nil {
			return err
		}
		fmt.Print(string(filtered))
		return nil
	}
}

func runConfigClear(cmd *cobra.Command, args []string) error {
	dir, err := homedir.Dir()
	if err != nil {
		return err
	}
	path := dir + "/.dmc.yaml"
	err = os.Remove(path)
	if err != nil {
		return err
	}
	good("Config cleared!")
	return nil
}

func checkValid(token string) bool {
	parsed, _ := jwt.Parse(token, nil)
	if parsed == nil {
		return false
	}
	return parsed.Claims.Valid() == nil
}

func login() (string, error) {
	fmt.Println("Enter dronemissioncontrol login details")
	emailPrompt := &promptui.Prompt{
		Label: "Email",
	}
	email, err := emailPrompt.Run()
	if err != nil {
		return "", err
	}

	passPrompt := &promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}
	pass, err := passPrompt.Run()
	if err != nil {
		return "", err
	}

	jsonStr := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, pass)
	req, err := http.NewRequest(
		"POST",
		"https://api.dronemissioncontrol.com/drone/gettoken",
		bytes.NewReader([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	token := strings.Trim(string(body), "\"\n")
	token = strings.TrimRight(token, `"`)
	if len(token) < 10 {
		return "", errors.New("invalid email or password")
	}
	good("Login successful!")
	return token, nil
}

func token() (string, error) {
	if token := viper.GetString("TOKEN"); token != "" && checkValid(token) {
		return token, nil
	}
	token, err := login()
	if err != nil {
		return "", nil
	}
	viper.Set("TOKEN", token)
	if err := writeConfig(); err != nil {
		return "", err
	}
	return token, nil
}

type dronesResponse struct {
	Drones []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"drones"`
}

func getDrones(token string) (*dronesResponse, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.dronemissioncontrol.com/user/drones/",
		nil,
	)
	req.Header.Set("Authorization", token)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var drones dronesResponse
	if err := json.NewDecoder(resp.Body).Decode(&drones); err != nil {
		return nil, errors.New("bad response from server")
	}
	return &drones, nil
}

func droneConfig(drones *dronesResponse) (id, password, url string, err error) {
	names := make([]string, len(drones.Drones))
	for i, drone := range drones.Drones {
		names[i] = drone.Name
	}
	selectPrompt := &promptui.Select{
		Label: "Select drone",
		Items: names,
	}
	idx, _, err := selectPrompt.Run()
	if err != nil {
		return
	}
	id = drones.Drones[idx].ID

	passPrompt := &promptui.Prompt{
		Label: "Drone verification key",
		Mask:  '*',
	}
	password, err = passPrompt.Run()

	urlPromt := &promptui.Prompt{
		Label:   "FCU url",
		Default: "udp://:14650@",
	}
	url, err = urlPromt.Run()

	return
}

func writeConfig() error {
	err := viper.WriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			dir, err := homedir.Dir()
			if err != nil {
				return err
			}
			return viper.WriteConfigAs(dir + "/.dmc.yaml")
		}
	}
	return err
}

func envList(keys ...string) (list []string) {
	for _, k := range keys {
		if v := viper.GetString(k); v != "" {
			list = append(list, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return
}

func init() {
	rootCmd.AddCommand(configCmd, loginCmd)
	configCmd.AddCommand(droneCmd, obcCmd, anipCmd, listCmd, clearCmd)
}
