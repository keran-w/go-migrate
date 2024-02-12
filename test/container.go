package main

import (
	"fmt"
	"github.com/keran-w/go-migrate/dockerapi"
	"os"
	"time"
)

func main() {
	args := os.Args
	containerName := "m1-num-printer-container"
	newContainerName := "m1-num-printer-container-new"

	if args[1] == "A" {
		fmt.Println("A")
		container, err := dockerapi.NewContainer(containerName, "m1-num-printer", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		container.Start()
		container.CopyOutput()
		time.Sleep(105 * time.Second)
		
	} else {
		fmt.Println("B")
		container, err := dockerapi.FindContainer(containerName)
		if err != nil {
			fmt.Println(err)
			return
		}
		newContainer := container.Migrate(newContainerName)
		newContainer.CopyOutput()
		container.Stop()
		time.Sleep(5 * time.Second)
		newContainer.Stop()
	}
}
