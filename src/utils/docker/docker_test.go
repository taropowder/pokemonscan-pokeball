package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetImages(t *testing.T) {
	//()
	CleanHangContainers()
}

func TestDockerRun(t *testing.T) {
	var err error
	exports := make(nat.PortSet, 10)
	port, err := nat.NewPort("tcp", "7777")
	if err != nil {
		log.Fatal(err)
	}
	exports[port] = struct{}{}

	containerConfig := &container.Config{
		Image: "pokemon:plugin_xray",
		Cmd: []string{"webscan", "--listen", "0.0.0.0:7777",
			"--webhook-output", fmt.Sprintf("http://host.docker.internal:%d/webhook", 5212)},
		WorkingDir:   "/app/workdir",
		ExposedPorts: exports,
		Hostname:     "pokemon-xray",
	}

	hostConfig := &container.HostConfig{ExtraHosts: []string{"host.docker.internal:host-gateway"},
		PortBindings: nat.PortMap{
			nat.Port("7777/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "7777",
				},
			},
		}, Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/Users/taro/GolandProjects/pokemon/data/passive_xray/config/config.yaml",
				Target: "/app/workdir",
			},
		}}

	Run(containerConfig, hostConfig, nil, "pokemon-xray")
}

func TestDockerWaitForRun(t *testing.T) {
	var err error
	exports := make(nat.PortSet, 10)
	port, err := nat.NewPort("tcp", "7777")
	if err != nil {
		log.Fatal(err)
	}
	exports[port] = struct{}{}

	containerConfig := &container.Config{
		Image: "pokemon:plugin_xray",
		Cmd: []string{"webscan", "--listen", "0.0.0.0:7777",
			"--webhook-output", fmt.Sprintf("http://host.docker.internal:%d/webhook", 5212)},
		WorkingDir:   "/app/workdir",
		ExposedPorts: exports,
		Hostname:     "pokemon-xray",
	}

	hostConfig := &container.HostConfig{ExtraHosts: []string{"host.docker.internal:host-gateway"},
		PortBindings: nat.PortMap{
			nat.Port("7777/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "7777",
				},
			},
		}}

	err = WaitForRun(containerConfig, hostConfig, nil, "pokemon-xray")
	if err != nil {
		log.Error(err)
	}
}

func TestDockerRm(t *testing.T) {
	RmWithContainerName("pokemon-xray")

}

func TestDockerExist(t *testing.T) {
	fmt.Println(ContainerExist("52b54fff44be"))

}
