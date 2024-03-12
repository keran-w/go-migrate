package main

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/keran-w/go-migrate/docker"
	// "github.com/keran-w/go-migrate/server"
	"time"
)

func main() {
	startTime := time.Now()

	imageName := "m1-number-printer-image:1.0"
	containerName := "m1-number-printer-container-B"
	env := []string{"START=0", "END=3000"}
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

	checkpointID := "checkpointA-1"
	checkpointDir := "/home/ubuntu/go-migrate/checkpoints"
	src := filepath.Join(checkpointDir, checkpointID)

	dst := filepath.Join("/var/lib/docker/containers", container.ID, "checkpoints", checkpointID)


	cmd := exec.Command("sudo", "mv", src, dst)
    err = cmd.Run()
    if err != nil {
        log.Fatalf("Error in transmitting checkpoints: %v", err)
        return
    }

	err = container.Restore(checkpointID, checkpointDir)
	if err != nil {
		log.Fatalf("Error restoreing from checkpoint %s: %v", checkpointID, err)
	}
	container.Start()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("Time taken for resuming: %v\n", duration)
}
