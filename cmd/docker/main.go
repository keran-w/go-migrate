package main

import (
	"github.com/keran-w/go-migrate/docker"
	"log"
)

func main() {

	imageName := "m1-number-printer-image:1.0"
	containerName := "m1-number-printer-container-A"
	env := []string{"START=50", "END=3000"}
	container, err := docker.NewContainer(containerName, imageName, env)
	if err != nil {
		log.Fatalf("Error creating container %s: %v", containerName, err)
		return
	}
	container.Start()
	// container.Stop()
}
