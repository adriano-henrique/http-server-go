package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	fmt.Println("Server listening on 0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()
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
		return
	}
	requestContent := string(readBuffer)
	requestHeader := ParseRequest(requestContent)
	fmt.Print("Received request: \n")
	requestHeader.prettyPrint()
	var response []byte
	path := requestHeader.method.url
	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(path, "/echo") {
		response = handleEcho(path)
	} else if strings.HasPrefix(path, "/user-agent") {
		response = handleUserAgent(requestHeader.userAgent)
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("found an error trying to respond")
		return
	}
}

func handleUserAgent(userAgent string) []byte {
	outputString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-length: %d\r\n\r\n%s\r\n", len(userAgent), userAgent)
	return []byte(outputString)
}

func handleEcho(url string) []byte {
	urlParams := strings.SplitN(url, "/", 3)
	echoedString := urlParams[2]
	outputString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-length: %d\r\n\r\n%s\r\n", len(echoedString), echoedString)

	fmt.Printf("Response:\n%s", outputString)

	return []byte(outputString)
}
