package server

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

func StartServer(netType, host, port string) {
	log.Println("Server Running...")
	server, err := net.Listen(netType, host+":"+port)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
		os.Exit(1)
	}
	defer server.Close()
	log.Println("Listening on " + host + ":" + port)

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %v", err)
			continue
		}
		log.Println("Client connected")
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	for i := 0; i < 3; i++ {
		var msg Message
		err := json.NewDecoder(connection).Decode(&msg)
		if err != nil {
			log.Fatalf("Error reading: %v", err)
			return
		}
		log.Printf("Received: %#v\n", msg)

		response := Message{Text: "Thanks! Got your message: " + msg.Text}
		err = json.NewEncoder(connection).Encode(response)
		if err != nil {
			log.Fatalf("Error sending response: %v", err)
			return
		}
	}
	log.Println("Finished 3 interactions, closing connection.")
}
