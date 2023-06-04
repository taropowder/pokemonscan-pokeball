package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const maxRunningSecond = 60 * 60

var dockerClient *client.Client

func init() {
	var err error
	//dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	dockerClient, err = client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		log.Fatalf("init docker error %s", err)
	} else {
		log.Info("docker client init success")
	}
}

func CleanHangContainers() {

	log.Infof("start to clean containers")
	filter := filters.NewArgs()

	ctx := context.Background()

	now := time.Now().Unix()
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
		log.Errorf("error GetRunningContainers %s", err)
	}

	for _, runningContainer := range containers {
		containerNames := strings.Join(runningContainer.Names, ",")
		if strings.Contains(containerNames, "pokemon-") && !strings.Contains(containerNames, "daemon") {
			runningTime := now - runningContainer.Created
			if runningTime > maxRunningSecond {

				removeOptions := types.ContainerRemoveOptions{
					RemoveVolumes: true,
					Force:         true,
				}

				if err := dockerClient.ContainerRemove(ctx, runningContainer.ID, removeOptions); err != nil {
					log.Errorf("Unable to remove container: %s", err)
				} else {
					log.Infof("clean %v", runningContainer.Names)
				}
			}
		}
	}
}

func Run(containerConfig *container.Config, containerHostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (err error) {
	ctx := context.Background()

	if networkingConfig == nil {
		networkingConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{"pokemon_net": {NetworkID: "pokemon_net"}}}
	}
	resp, err := dockerClient.ContainerCreate(ctx, containerConfig, containerHostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return err
	}
	log.Infof("container %s create ", containerName)
	err = dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err == nil {
		log.Infof("container %s start: %s ", containerName, resp.ID)
	}
	return
}

func WaitForRun(containerConfig *container.Config, containerHostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (err error) {
	ctx := context.Background()

	if networkingConfig == nil {
		networkingConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{"pokemon_net": {NetworkID: "pokemon_net"}}}
	}
	resp, err := dockerClient.ContainerCreate(ctx, containerConfig, containerHostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return err
	}
	log.Infof("container %s create ", containerName)
	err = dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err == nil {
		log.Infof("container %s start: %s ", containerName, resp.ID)
	}
	statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Errorf("error for run %s:%v", containerName, err)
		}
	case <-statusCh:
	}
	log.Infof("container %s stop ", containerName)

	return
}

func CreateDockerConfig(pluginName, forwardPort, cwd string, mounts []mount.Mount,
	command []string) (containerConfig *container.Config,
	containerHostConfig *container.HostConfig) {

	//var err error
	exports := make(nat.PortSet, 10)
	if forwardPort != "" {
		port, err := nat.NewPort("tcp", forwardPort)
		if err != nil {
			log.Fatal(err)
		}
		exports[port] = struct{}{}
	}

	containerConfig = &container.Config{
		Image:        fmt.Sprintf("pokemon:plugin_%s", pluginName),
		Cmd:          command,
		WorkingDir:   cwd,
		ExposedPorts: exports,
		Hostname:     fmt.Sprintf("pokemon:plugin_%s", pluginName),
	}

	//portBindings := nat.PortMap{}
	//if forwardPort != "" {
	//	portBindings = nat.PortMap{
	//		nat.Port("7777/tcp"): []nat.PortBinding{
	//			{
	//				HostIP:   "0.0.0.0",
	//				HostPort: "7777",
	//			},
	//		},
	//	}
	//}

	//hostConfig := &container.HostConfig{ExtraHosts: []string{"host.docker.internal:host-gateway"},
	//	PortBindings: portBindings,
	//	Mounts: mounts}

	return
}

func RmWithContainerName(containerName string) error {
	ctx := context.Background()

	err := dockerClient.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		Force: true,
	})

	return err
}

func ContainerExist(containerName string) bool {
	ctx := context.Background()

	i, err := dockerClient.ContainerInspect(ctx, containerName)

	// 判断容器是否运行， 如果是 running 状态的容器则返回 true
	log.Info(i.State.Status)

	if err == nil && i.State.Status == "running" {
		return true
	}

	return false
}
