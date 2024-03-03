package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// newDockerClient creates a new Docker client from environment variables.
func newDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// createContainer creates a new container from the specified image and environment variables.
func createContainer(ctx context.Context, cli *client.Client, containerName, imageName string, env []string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
	}, nil, nil, nil, containerName)
	return resp, err
}

func inspectContainer(cli *client.Client, containerName string) (types.ContainerJSON, error) {
	return cli.ContainerInspect(context.Background(), containerName)
}

// deleteContainer removes a container by its ID.
func deleteContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

// startContainer starts a container by its ID.
func startContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

// stopContainer stops a container by its ID.
func stopContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStop(ctx, containerID, container.StopOptions{})
}
