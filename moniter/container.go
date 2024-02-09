package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	// "github.com/docker/docker/pkg/stdcopy"
	// "os"
	// "time"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containerName := "m1-num-printer-container"
	imageName := "m1-number-printer"

	// containers, err := cli.ContainerList(ctx, container.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// var containerID string
	// for _, container := range containers {
	// 	for _, name := range container.Names {
	// 		if name == "/"+containerName {
	// 			containerID = container.ID
	// 			break
	// 		}
	// 	}
	// }

	// if containerID == "" {
	// 	fmt.Printf("Container named %s not found\n", containerName)
	// 	return
	// }

	// processes, err := cli.ContainerTop(ctx, containerID, nil)
	// if err != nil {
	// 	fmt.Println("Error getting container processes:", err)
	// 	return
	// }

	// fmt.Println("Processes running in the container:")
	// for _, process := range processes.Processes {
	// 	for j, col := range process {
	// 		fmt.Printf("%s: %s ", processes.Titles[j], col)
	// 	}
	// 	fmt.Println()
	// }

	// timeout := 1
	// stopOptions := container.StopOptions{
	// 	Timeout: &timeout,
	// }

	// if err := cli.ContainerStop(ctx, containerID, stopOptions); err != nil {
	// 	fmt.Println("Error stopping container:", err)
	// 	return
	// }

	// fmt.Printf("Container %s stopped\n", containerID)

	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true}); err != nil {
		fmt.Println("Container does not exist", err)
	}

	containerConfig := &container.Config{
		Image: imageName,
		Env:   []string{"START=1"},
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, containerName)
	if err != nil {
		fmt.Println("Error creating container:", err)
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Error starting container:", err)
		return
	}

	fmt.Printf("Container %s started with START=50\n", containerName)
}
