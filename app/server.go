package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	fmt.Println("Server listening on 0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	readBuffer := make([]byte, 1024)
	_, err := conn.Read(readBuffer)
	if err != nil {
		fmt.Println("error reading buffer: ", err.Error())
		os.Exit(1)
	}
	fmt.Print("Responding with HTTP/1.1 200 OK\r\n\r\n")
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
