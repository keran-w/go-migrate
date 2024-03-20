package main

import (
	"github.com/keran-w/go-migrate/docker"
	"log"
	"os/exec"
	"time"
)

func main() {

	containerName := "m1-number-printer-container-A"
	container, err := docker.FindContainer(containerName)
	if err != nil {
		log.Fatalf("Error finding container %s: %v", containerName, err)
		return
	}

	log.Printf("Creating checkpoint for container %s...\n", containerName)
	// checkpointName := "checkpointA-" + time.Now().Format("MM-DDTHH-mm")
	checkpointName := "checkpointA-1"
	checkpointDir := "/home/ubuntu/go-migrate/checkpoints"

	startTime := time.Now()

	err = container.Checkpoint(checkpointName, checkpointDir, false)
	if err != nil {
		log.Fatalf("Error creating checkpoint for container %s: %v", containerName, err)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("Time taken for checkpoint: %v\n", duration)

	cmd := exec.Command("sudo", "chmod", "-R", "777", "./checkpoints")
	if cmd.Run() != nil {
		log.Fatalf("Error changing permissions for checkpoint directory: %v", err)
		return
	} else {
		log.Printf("Permissions changed for checkpoint directory.\n")
	}
	


	time.Sleep(5 * time.Second)


	log.Printf("Creating checkpoint for container %s...\n", containerName)
	// checkpointName := "checkpointA-" + time.Now().Format("MM-DDTHH-mm")
	checkpointName = "checkpointA-2"
	checkpointDir = "/home/ubuntu/go-migrate/checkpoints"

	startTime = time.Now()

	err = container.Checkpoint(checkpointName, checkpointDir, false)
	if err != nil {
		log.Fatalf("Error creating checkpoint for container %s: %v", containerName, err)
		return
	}

	endTime = time.Now()
	duration = endTime.Sub(startTime)
	log.Printf("Time taken for checkpoint: %v\n", duration)

	cmd = exec.Command("sudo", "chmod", "-R", "777", "./checkpoints")
	if cmd.Run() != nil {
		log.Fatalf("Error changing permissions for checkpoint directory: %v", err)
		return
	} else {
		//log.Printf("Permissions changed for checkpoint directory.\n")
	}

	// varName := "CURR"
	// value := container.GetState(varName)
	// log.Printf("Container %s state %s: %s\n", containerName, varName, value)

	//netType := "tcp"
	//host := "localhost"
	//port := "9988"
	//client.ConnectToServer(netType, host, port)
}
