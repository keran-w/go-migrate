package main

import (
	"github.com/keran-w/go-migrate/docker"
	"log"
)

func main() {

	imageName := "ml_app"
	containerName := "ml-container-A"
	env := []string{"START=0", "END=3000"}
	container, err := docker.NewContainer(containerName, imageName, env)
	if err != nil {
		log.Fatalf("Error creating container %s: %v", containerName, err)
		return
	}
	container.Start()
	// container.Stop()
}
