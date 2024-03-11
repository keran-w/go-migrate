package client

// build image (sh scripts/build_images) 
// -> run docker/main.go (new container running the image)
// client create checkpoint -> send checkpoint -> server receive checkpoint -> process
// -> server evaluate -> server sends evaluate results to client
// if finish, client stops the container and send remaining pages
// if not, repeat

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Message struct {
	Text string `json:"text"`
}

func ConnectToServer(netType, host, port string) {
	connection, err := net.Dial(netType, host+":"+port)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer connection.Close()

	for i := 0; i < 3; i++ {
		msg := Message{Text: fmt.Sprintf("Hello Server! Greetings #%d.", i+1)}
		err := json.NewEncoder(connection).Encode(msg)
		if err != nil {
			log.Fatalf("Error sending message to server: %v", err)
			break
		}

		var response Message
		err = json.NewDecoder(connection).Decode(&response)
		if err != nil {
			log.Fatalf("Error reading response from server: %v", err)
			break
		}

		log.Printf("Received from server: %#v", response)
	}
}
