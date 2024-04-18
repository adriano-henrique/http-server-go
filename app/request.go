package main

import (
	"fmt"
	"strings"
)

type HttpMethod struct {
	verb     string
	url      string
	protocol string
}

type RequestHeader struct {
	method    HttpMethod
	host      string
	userAgent string
}

func removeWhitespaceFromEOF(list []string) []string {
	var result []string
	for i := 0; i < len(list); i++ {
		if strings.TrimSpace(list[i]) != "" {
			result = append(result, list[i])
		}
	}
	return result
}

func ParseRequest(requestContent string) *RequestHeader {
	headerEndIndex := strings.Index(requestContent, "\r\n\r\n")
	var header string
	if headerEndIndex != -1 {
		header = requestContent[:headerEndIndex]
	} else {
		// If double newline is not found, consider the whole string as header
		header = requestContent
	}
	parsedRequest := removeWhitespaceFromEOF(strings.Split(header, "\r\n"))
	methodLine := strings.Split(parsedRequest[0], " ")

	httpMethod := HttpMethod{
		verb:     strings.TrimSpace(methodLine[0]),
		url:      strings.TrimSpace(methodLine[1]),
		protocol: strings.TrimSpace(methodLine[2]),
	}

	var host string
	var userAgent string
	if len(parsedRequest) > 1 {
		fmt.Println("Im here")
		host = strings.Split(parsedRequest[1], ": ")[1]
		userAgent = strings.Split(parsedRequest[2], ": ")[1]
	}
	fmt.Println(host)
	fmt.Println(userAgent)

	return &RequestHeader{
		method:    httpMethod,
		host:      host,
		userAgent: userAgent,
	}
}

func (r *RequestHeader) prettyPrint() {
	fmt.Printf("Method: %s %s %s\n", r.method.verb, r.method.url, r.method.protocol)
	fmt.Printf("Host: %s\n", r.host)
	fmt.Printf("User-Agent: %s\n", r.userAgent)
}
