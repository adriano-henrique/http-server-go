package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	dirFlag := flag.String("directory", "", "The directory to call /files http")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	fmt.Println("Server listening on 0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()
	flag.Parse()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, *dirFlag)
	}
}

func handleConnection(conn net.Conn, directory string) {
	fmt.Println(directory)
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
	} else if strings.HasPrefix(path, "/files") {
		response = handleGetFileContents(path, directory)
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

func handleGetFileContents(urlParams string, directory string) []byte {
	file := strings.Split(urlParams, "/")[2]
	if directory == "" {
		return []byte(buildErrorResponse("Should pass directory to command"))
	}

	dir, err := os.Open(directory)
	if err != nil {
		fmt.Println(err)
		return []byte(buildErrorResponse(err.Error()))
	}
	files, err := dir.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return []byte(buildErrorResponse(err.Error()))
	}
	var hasMatch bool
	for _, v := range files {
		if !v.IsDir() && v.Name() == file {
			hasMatch = true
		}
	}
	if !hasMatch {
		return []byte(buildErrorResponse("File is not on directory"))
	}

	fileContent, err := os.ReadFile(directory + file)
	if err != nil {
		return []byte(buildErrorResponse(err.Error()))
	}
	outputString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-length: %d\r\n\r\n%s\r\n", len(fileContent), fileContent)
	return []byte(outputString)
}

func buildErrorResponse(data string) string {
	return fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n%s", data)
}
