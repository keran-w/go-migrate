package dockerapi

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"os"
	"strings"
	// "time"
)

// UNUSED is a utility function to explicitly ignore unused variables.
func UNUSED(x ...interface{}) {}

// Container represents a Docker container with basic configuration.
type Container struct {
	Name  string
	Image string
	Env   []string
}

// newDockerClient initializes and returns a new Docker client.
// Returns the Docker client and any error encountered.
func newDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return cli, err
}

// FindContainer is a helper function that searches for a container by name.
// It returns the found Container struct and nil if the container exists, or an empty Container struct and an error if not.
func FindContainer(containerName string) (Container, error) {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		return Container{}, err
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return Container{}, err
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName {
				containerDetails, err := cli.ContainerInspect(ctx, container.ID)
				if err != nil {
					return Container{}, err
				}
				return Container{
					Name:  containerName,
					Image: containerDetails.Config.Image,
					Env:   containerDetails.Config.Env,
				}, nil
			}
		}
	}

	return Container{}, fmt.Errorf("container named %s not found", containerName)
}

// NewContainer creates a new Container instance by either finding an existing container or creating a new one.
// It returns a pointer to the Container struct and an error if any.
func NewContainer(containerName, imageName string, env []string) (*Container, error) {
	// Attempt to find the container first
	foundContainer, err := FindContainer(containerName)
	if err == nil {
		fmt.Println("Container found, no need to create a new one.")
		return &foundContainer, nil
	}

	// If not found, create a new container
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		return nil, err
	}

	// Create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
	}, nil, nil, nil, containerName)
	if err != nil {
		return nil, err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, err
	}

	fmt.Printf("New container %s created and started.\n", containerName)
	return &Container{Name: containerName, Image: imageName, Env: env}, nil
}

// containerStart starts a container by its name. If it does not exist, it creates and starts a new container with the given image and environment variables.
func (c *Container) Start() {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return
	}

	if err := cli.ContainerStart(ctx, c.Name, container.StartOptions{}); err != nil {
		fmt.Println("Error starting container:", err)
		return
	}
	fmt.Printf("Container %s started\n", c.Name)
}

// stop stops and removes the specified container.
// Logs any errors encountered during the stop or remove operations.
func (c *Container) Stop() {
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return
	}

	if err := cli.ContainerStop(ctx, c.Name, container.StopOptions{}); err != nil {
		fmt.Println("Error stopping container:", err)
		return
	}

	fmt.Printf("Container %s stopped\n", c.Name)
}

// replaceOrAppendEnv searches for an environment variable in a slice of strings.
// If found, it replaces its value. If not, it appends the variable with the new value.
func replaceOrAppendEnv(env []string, key, newValue string) []string {
	keyPrefix := key + "="
	found := false

	for i, v := range env {
		if strings.HasPrefix(v, keyPrefix) {
			env[i] = keyPrefix + newValue
			found = true
			break
		}
	}

	if !found {
		env = append(env, keyPrefix+newValue)
	}

	return env
}

func (c *Container) Migrate(newContainerName string) *Container {
	fmt.Printf("Migrating container %s to %s\n", c.Name, newContainerName)
	currImage := c.Image

	updatedEnv := &container.Config{
        Image: currImage,
        Env:   []string{"START=50"},
    }
	// fmt.Printf("Updated Env: %v\n", currEnv)

	newContainer, err := NewContainer(newContainerName, currImage, updatedEnv.Env)
	if err != nil {
		fmt.Println("Error creating new container:", err)
		return nil
	}

	newContainer.Start()
	return newContainer
}

// copyOutput attaches to the container's output streams and copies them to the local standard output and standard error.
// It runs in a separate goroutine to not block the main execution flow.
func (c *Container) CopyOutput() {
	go func() {
		ctx := context.Background()
		cli, err := newDockerClient()
		if err != nil {
			fmt.Println("Error creating Docker client:", err)
			return
		}

		out, err := cli.ContainerAttach(ctx, c.Name, container.AttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		})
		if err != nil {
			fmt.Println("Error attaching to container:", err)
			return
		}
		defer out.Close()

		fmt.Printf("Attached to container %s\n", c.Name)

		if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out.Reader); err != nil {
			fmt.Println("Error copying container output:", err)
		}
	}()
}
