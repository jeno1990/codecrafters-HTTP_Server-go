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

	path_list := strings.Split(req.Path, "/")
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
	if req.Path == "/" || path_list[1] == "echo" || path_list[1] == "user-agent" {
		body, ok := req.Headers["User-Agent"]
		if !ok {
			body = path_list[len(path_list)-1]
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
	response.WriteResponse(conn, req)

}

// func response(conn net.Conn, okay bool, body string) {
// 	defer conn.Close()
// 	var response string = ""
// 	if okay {
// 		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
// 	} else {
// 		response = "HTTP/1.1 404 Not Found\r\n\r\n"
// 	}
// 	_, err := conn.Write([]byte(response))
// 	if err != nil {
// 		fmt.Println("Error writing to connection: ", err.Error())
// 		os.Exit(1)
// 	}
// }
