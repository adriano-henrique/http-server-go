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

func ParseRequest(requestContent string) *RequestHeader {
	parsedRequest := strings.Split(requestContent, "\r\n")
	methodLine := strings.Split(parsedRequest[0], " ")

	httpMethod := HttpMethod{
		verb:     strings.TrimSpace(methodLine[0]),
		url:      strings.TrimSpace(methodLine[1]),
		protocol: strings.TrimSpace(methodLine[2]),
	}

	return &RequestHeader{
		method:    httpMethod,
		host:      strings.Split(parsedRequest[1], ": ")[1],
		userAgent: strings.Split(parsedRequest[2], ": ")[1],
	}
}

func (r *RequestHeader) prettyPrint() {
	fmt.Printf("Method: %s %s %s\n", r.method.verb, r.method.url, r.method.protocol)
	fmt.Printf("Host: %s\n", r.host)
	fmt.Printf("User-Agent: %s\n", r.userAgent)
}
