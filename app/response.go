package main

import (
	"fmt"
	"net"
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
}

func (r *Response) WriteResponse(conn net.Conn, req *HttpRequest) {
	defer conn.Close()

	// Build the response headers
	headers := ""
	for key, value := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	content_length := ""
	if r.Body != "" {
		content_length = fmt.Sprintf("Content-Length: %d", len(r.Body))
	}

	// Build the full response
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n%s%s\r\n\r\n%s",
		r.StatusCode, r.StatusText, headers, content_length, r.Body)
	fmt.Print("response: ", response)
	// Write the response to the connection
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}
