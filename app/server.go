package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const ()

func main() {
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
	readValue, err := conn.Read(readBuffer)
	if err != nil {
		fmt.Println("error reading buffer: ", err.Error())
		os.Exit(1)
	}
	requestContent := string(readBuffer[:readValue])
	fmt.Print("Received request: ", requestContent)
	_, url, _ := parseStartLine(requestContent)
	if url == "/" {
		fmt.Print("Responding with HTTP/1.1 200 OK\r\n\r\n")
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(url, "/echo") {
		handleEcho(conn, url)
	} else {
		fmt.Print("Responding with HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func parseStartLine(requestContent string) (string, string, string) {
	startLine := strings.Split(requestContent, "\r\n")[0]
	requestParts := strings.Split(startLine, " ")
	return strings.TrimSpace(requestParts[0]), strings.TrimSpace(requestParts[1]), strings.TrimSpace(requestParts[2])
}

func handleEcho(conn net.Conn, url string) {
	urlParams := strings.SplitN(url, "/", 3)
	echoedString := urlParams[2]
	outputString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-length: %d\r\n\r\n%s\r\n", len(echoedString), echoedString)

	fmt.Printf("Response:\n%s", outputString)

	conn.Write([]byte(outputString))
}
