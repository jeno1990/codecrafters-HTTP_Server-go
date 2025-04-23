package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
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
		go request(conn)
	}

}

func request(conn net.Conn) {
	defer conn.Close()
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	str := string(request[:n])
	req, err := parseRequestLine(str)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}

	url := strings.Split(req.Path, "/")
	fmt.Println("req: ", req)

	// if req.Path == "/" || path_list[1] == "echo" || path_list[1] == "user-agent" {
	// 	userAgent, ok := req.Headers["User-Agent"]
	// 	if !ok {
	// 		userAgent = path_list[len(path_list)-1]
	// 	}
	// 	response(conn, true, userAgent)
	// } else {
	// 	response(conn, false, "")
	// }
	var response Response
	if req.Method == "GET" {

		if url[1] == "files" {
			fileName := url[2]
			if _, err := os.Stat(os.Args[2] + fileName); err == nil {
				// File exists, read its contents
				content, _ := os.ReadFile(os.Args[2] + fileName)
				response = Response{
					StatusCode: 200,
					StatusText: "OK",
					Headers: map[string]string{
						"Content-Type": "application/octet-stream",
					},
					Body: string(content),
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
			body, ok := req.Headers["User-Agent"]
			if !ok {
				body = url[len(url)-1]
			}
			response = Response{
				StatusCode: 200,
				StatusText: "OK",
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: body,
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
		print("POST request------------------------")
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
	response.WriteResponse(conn, req)
}
