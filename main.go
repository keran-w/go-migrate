package main

import (
	"github.com/checkpoint-restore/go-criu/v7"
	"log"
)

func main() {
	c := criu.MakeCriu()
	version, err := c.GetCriuVersion()
	if err != nil {
		log.Println("Error getting CRIU version:")
		log.Fatalln(err)
	} else {
		log.Println("CRIU version:")
		log.Println(version)
	}
}
