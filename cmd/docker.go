package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

var (
	dockerClient *client.Client
)

const (
	imageBase = "docker.io/tobiasfriden/"
)

const dockerFailMessage = `
Could not connect to docker. Docker is necessary in order to run Drone Mission Control onboard software.

If docker is not installed, please install it by following the instructions at: https://docs.docker.com/install/

If docker is running, you might need to run this CLI as sudo. You can also add permissions for the current user to use docker:

  sudo usermod -aG docker USER

Note that these changes require logout to take affect.
`

func init() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Print(dockerFailMessage)
		os.Exit(1)
	}
	dockerClient = cli
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

func pullImage(name, imageName string) error {
	fmt.Printf("Pulling %s..\n", name)
	out, err := dockerClient.ImagePull(context.Background(), imageBase+imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	if Verbose {
		err := jsonmessage.DisplayJSONMessagesStream(
			out,
			os.Stdout,
			os.Stdout.Fd(),
			true,
			func(message jsonmessage.JSONMessage) {},
		)
		if err != nil {
			return err
		}
	} else {
		_, err = io.Copy(ioutil.Discard, out)
		if err != nil {
			return err
		}
	}
	good("Done!")
	return err
}

func startContainer(name, imageName string, config *container.Config, hostConfig *container.HostConfig) error {
	running, err := containerRunning(imageName)
	if err != nil {
		return err
	}
	if running {
		fmt.Printf("Container %s is already running\n", name)
		if !Recreate {
			return nil
		} else {
			if err := stopContainer(name, imageName); err != nil {
				return err
			}
		}
	}
	fmt.Printf("Creating %s..\n", name)
	config.Image = imageBase + imageName
	ctx := context.Background()
	resp, err := dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		nil,
		name,
	)
	if err != nil {
		return err
	}
	fmt.Printf("Starting %s..\n", name)
	if err := dockerClient.ContainerStart(
		ctx,
		resp.ID,
		types.ContainerStartOptions{},
	); err != nil {
		return err
	}
	good("Done!")
	return nil
}

func containerRunning(imageName string) (bool, error) {
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return false, err
	}
	for _, c := range containers {
		if c.Image == imageBase+imageName {
			return true, nil
		}
	}
	return false, nil
}

func stopContainer(name, imageName string) error {
	fmt.Printf("Stopping %s..\n", name)
	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	for _, c := range containers {
		if c.Image == imageBase+imageName {
			if err := dockerClient.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				return err
			}
		}
	}
	good("Done!")
	return nil
}
