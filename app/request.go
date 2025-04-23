package main

import (
	"log"
	"strings"
)

type HttpRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func (req *HttpRequest) GetHeader(key string) string {
	if value, ok := req.Headers[key]; ok {
		return value
	}
	return ""
}

func parseRequestLine(requestLine string) (*HttpRequest, error) {
	req := HttpRequest{}
	parts := strings.Split(requestLine, "\r\n")
	// print("reqLine: ===> ", requestLine)
	method_path_version := strings.Fields(parts[0])
	host := strings.Fields(parts[1])
	headers := make(map[string]string)
	headers["Host"] = host[1]

	for i := len(parts) - 2; i > 0; i-- {
		header := strings.Split(parts[i], ":")
		if len(header) == 2 {
			headers[header[0]] = strings.TrimSpace(header[1])
		}
	}

	req.Method = method_path_version[0]
	req.Path = method_path_version[1]
	req.Version = method_path_version[2]
	req.Headers = headers
	req.Body = parts[len(parts)-1]
	log.Println("req: ", req)
	return &req, nil
}
