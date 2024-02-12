package main

import (
	"fmt"
	"github.com/keran-w/go-migrate/dockerapi"
)

func main() {
	containerName := "m1-num-printer-container"
	container, err := dockerapi.NewContainer(containerName, "m1-num-printer", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	container.Start()
	container.CopyOutput()
	// time.Sleep(1 * time.Second)
	container.Stop()
}
