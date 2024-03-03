package main

import (
	"github.com/keran-w/go-migrate/docker"
	"log"
)

func main() {
	//netType := "tcp"
	//host := "localhost"
	//port := "9988"
	//client.ConnectToServer(netType, host, port)

	containerName := "m1-number-printer-container-A"
	container, err := docker.FindContainer(containerName)
	if err != nil {
		log.Fatalf("Error finding container %s: %v", containerName, err)
		return
	}
	varName := "CURR"
	value := container.GetState(varName)
	log.Printf("Container %s state %s: %s\n", containerName, varName, value)
}