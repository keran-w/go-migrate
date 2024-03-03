package docker

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/keran-w/go-migrate/utils"
	"log"
	"os"
	"strings"
)

// Container represents a Docker container with basic configuration.
type Container struct {
	Name  string
	Image string
	Env   []string
}

// NewContainer ensures a container exists and is started, either by finding it or creating a new one.
func NewContainer(containerName, imageName string, env []string) (*Container, error) {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	ctr, err := FindContainer(containerName)
	if err == nil && ctr != nil {
		log.Println("Container found, removing it...")
		err := deleteContainer(ctx, cli, containerName)
		if err != nil {
			log.Fatalf("Error deleting container %s: %v", containerName, err)
			return nil, err
		}
		// Reuse the existing container
		// return &Container{Name: containerName, Image: containerDetails.Config.Image, Env: containerDetails.Config.Env}, nil
	}

	resp, err := createContainer(ctx, cli, containerName, imageName, env)
	if err != nil {
		log.Fatalf("Error creating container %s: %v", containerName, err)
		return nil, err
	}

	utils.UNUSED(resp)

	log.Printf("New container %s created.\n", containerName)
	return &Container{Name: containerName, Image: imageName, Env: env}, nil
}

func (c *Container) Start() {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	if err := startContainer(ctx, cli, c.Name); err != nil {
		log.Fatalf("Error starting container %s: %v", c.Name, err)
	}

	log.Printf("Container %s started.\n", c.Name)
}

func (c *Container) Stop() {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	if err := stopContainer(ctx, cli, c.Name); err != nil {
		log.Fatalf("Error stopping container %s: %v", c.Name, err)
	}

	log.Printf("Container %s stopped.\n", c.Name)
}

func (c *Container) GetState(varName string) string {
	// TODO: GetState returns the initial state of the container, not the current state.
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	containerInfo, err := inspectContainer(cli, c.Name)
	if err != nil {
		log.Fatalf("Error inspecting container %s: %v", c.Name, err)
	}
	envVariables := containerInfo.Config.Env
	for _, envVar := range envVariables {
		envVarName := strings.Split(envVar, "=")[0]
		if envVarName == varName {
			return strings.Split(envVar, "=")[1]
		}
	}
	return ""
}

func (c *Container) CopyOutput() {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	out, err := cli.ContainerAttach(ctx, c.Name, container.AttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		log.Fatalf("Error attaching to container: %v", err)
	}
	defer out.Close()

	log.Printf("Attached to container %s output.\n", c.Name)

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out.Reader)
	if err != nil {
		log.Printf("Error copying container output: %v", err)
	}
}

// FindContainer searches for a container by name.
func FindContainer(containerName string) (*Container, error) {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, ctr := range containers {
		for _, name := range ctr.Names {
			if name == "/"+containerName {
				containerDetails, err := cli.ContainerInspect(ctx, ctr.ID)
				if err != nil {
					return nil, err
				}
				return &Container{Name: containerName, Image: containerDetails.Config.Image, Env: containerDetails.Config.Env}, nil
			}
		}
	}
	return nil, nil
}
