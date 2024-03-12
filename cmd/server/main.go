package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/keran-w/go-migrate/docker"
	// "github.com/keran-w/go-migrate/server"
)

func main() {
	act := os.Args[1]
	switch act {
	case "restore":
		imageName := "m1-number-printer-image:1.0"
		containerName := "m1-number-printer-container-B"
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

		// checkpointID := "checkpoint"
		// checkpointDir := "/home/ubuntu/test/server"
		// src := filepath.Join(checkpointDir, checkpointID)
		checkpointDir := os.Args[2]
		checkpointID := os.Args[3]
		src := filepath.Join(checkpointDir, checkpointID)

		signalFilename := "signal.txt"
		signalDir := "/home/ubuntu/test/server"
		signalPath := filepath.Join(signalDir, signalFilename)

		for {
			fmt.Println("Server Waiting for Migration...")
			_, err := os.Stat(signalPath)
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
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
		container.Start()
	}
}
