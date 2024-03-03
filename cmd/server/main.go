package main

import "github.com/keran-w/go-migrate/server"

func main() {
	netType := "tcp"
	host := "localhost"
	port := "9988"
	server.StartServer(netType, host, port)
}
