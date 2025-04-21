package main

import (
	"fmt"
	"strings"
)

type HttpRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	// Body    string
}

func (req *HttpRequest) GetHeader(key string) string {
	if value, ok := req.Headers[key]; ok {
		return value
	}
	return ""
}

func parseRequestLine(requestLine string) (string, string, string, error) {
	fmt.Println("requestLine: -> ", requestLine)
	parts := strings.Split(requestLine, "\r\n")[0]
	fmt.Println("parts: -> ", parts)
	requestParts := strings.Fields(parts)
	if len(requestParts) < 3 {
		return "", "", "", fmt.Errorf("Invalid request line")
	}
	fmt.Println("requests: -> ", requestParts)
	method := requestParts[0]
	path := requestParts[1]
	version := requestParts[2]
	return method, path, version, nil
}
