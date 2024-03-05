package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/keran-w/go-migrate/docker"
	// "github.com/keran-w/go-migrate/server"
)

func main() {
	imageName := "m1-number-printer-image:1.0"
	containerName := "m1-number-printer-container-B-clone"
	env := []string{"START=50", "END=3000"}
	container, err := docker.NewContainer(containerName, imageName, env)
	if err != nil {
		log.Fatalf("Error creating container %s: %v", containerName, err)
		return
	}

	// TODO: communications
	// netType := "tcp"
	// host := "localhost"
	// port := "9988"
	// server.StartServer(netType, host, port)

	checkpointDir := "/tmp"
	checkpointID := "checkpointB"
	src := filepath.Join(checkpointDir, checkpointID)
	dst := filepath.Join("/var/lib/docker/containers", container.ID, "checkpoints", checkpointID)
	err = os.Rename(src, dst)
	if err != nil {
		log.Fatalf("Error in transmitting checkpoints: %v", err)
		return
	}

	err = container.Restore(checkpointID, checkpointDir)
	if err != nil {
		log.Fatalf("Error restoreing from checkpoint %s: %v", checkpointID, err)
	}
}
