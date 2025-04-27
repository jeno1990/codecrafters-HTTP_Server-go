package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func main() {

	listner, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer listner.Close()
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		request := make([]byte, 1024)
		n, err := conn.Read(request)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("client closed connection")
			} else if ne, ok := err.(net.Error); ok && ne.Timeout() {
				fmt.Println("read timeout; closing connection")
			} else {
				fmt.Printf("bad request: %v", err)
			}
			return
		}
		str := string(request[:n])
		req, err := parseRequestLine(str)
		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			return
		}
		response := processRequest(req)
		// fmt.Println("response: ", response)

		response.WriteResponse(conn, req)

		// break the loop if the request header contains the "Connection: close" header
		if req.GetHeader("Connection") == "close" {
			fmt.Println("Connection closed by client request")
			return
		}
	}
	// conn.Close()
}

func processRequest(req *HttpRequest) Response {
	url := strings.Split(req.Path, "/")
	// fmt.Println("req: ", req)

	var response Response
	headers := map[string]string{}
	if _, ok := req.Headers["Accept-Encoding"]; ok {
		encoddings := strings.Split(req.Headers["Accept-Encoding"], ",")
		for _, encoding := range encoddings {
			if strings.TrimSpace(encoding) == "gzip" {
				headers["Content-Encoding"] = "gzip"
			}
		}
	}
	if req.GetHeader("Connection") == "close" {
		headers["Connection"] = "close"
	}
	if req.Method == "GET" {
		if url[1] == "files" {
			headers["Content-Type"] = "application/octet-stream"
			fileName := url[2]
			if _, err := os.Stat(os.Args[2] + fileName); err == nil {
				// File exists, read its contents
				content, _ := os.ReadFile(os.Args[2] + fileName)
				response = Response{
					StatusCode: 200,
					StatusText: "OK",
					Headers:    headers,
					Body:       string(content),
				}

			} else if os.IsNotExist(err) {
				response = Response{
					StatusCode: 404,
					StatusText: "Not Found",
					Headers:    map[string]string{},
					Body:       "",
				}
			}

		} else if req.Path == "/" || url[1] == "echo" || url[1] == "user-agent" {
			body := ""
			if url[1] == "user-agent" {
				body, _ = req.Headers["User-Agent"]
			} else if url[1] == "echo" {
				body = url[len(url)-1]
			}
			// if !ok {
			// if url[1] == "echo" {
			// 	body = url[len(url)-1]
			// }
			headers["Content-Type"] = "text/plain"
			response = Response{
				StatusCode: 200,
				StatusText: "OK",
				Headers:    headers,
				Body:       body,
			}
		} else {
			response = Response{
				StatusCode: 404,
				StatusText: "Not Found",
				Headers:    map[string]string{},
				Body:       "",
			}
		}
	} else if req.Method == "POST" {
		file_path := os.Args[2] + url[2]
		os.WriteFile(file_path, []byte(req.Body), 0644)
		response = Response{
			StatusCode: 201,
			StatusText: "Created",
			Headers:    map[string]string{},
			Body:       "",
		}
	} else {
		response = Response{
			StatusCode: 405,
			StatusText: "Method Not Allowed",
			Headers:    map[string]string{},
			Body:       "",
		}
	}
	return response
}
