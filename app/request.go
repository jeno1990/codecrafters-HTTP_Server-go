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

	method_path_version := strings.Fields(parts[0])
	host := strings.Fields(parts[1])
	user_agent := strings.Fields(parts[2])
	origin := strings.Fields(parts[3])
	fmt.Println(method_path_version, " | ", host, " | ", origin, " | ", user_agent)
	req.Method = method_path_version[0]
	req.Path = method_path_version[1]
	req.Version = method_path_version[2]
	headers := make(map[string]string)
	headers["Host"] = host[1]
	if len(origin) > 1 {
		headers["Origin"] = origin[1]
	}
	if len(user_agent) > 1 {
		headers["User-Agent"] = user_agent[1]
	}
	req.Headers = headers
	return &req, nil
}
