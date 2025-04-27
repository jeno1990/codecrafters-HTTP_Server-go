package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
)

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
}

func (r *Response) WriteResponse(conn net.Conn, req *HttpRequest) {

	// Build the response headers
	headers := ""
	for key, value := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	content_length := ""
	var body = ""

	// compress the body if it exists
	if r.Body != "" {
		b, size, _ := r.compressBody(req)
		body = b
		content_length = fmt.Sprintf("Content-Length: %d", size)
	}

	// Build the full response
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n%s%s\r\n\r\n%s",
		r.StatusCode,
		r.StatusText,
		headers,
		content_length,
		body)
	fmt.Println("___response____: ", response)

	// Write the response to the connection
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}

func (r *Response) compressBody(req *HttpRequest) (string, int, error) {
	_, ok := req.Headers["Accept-Encoding"]
	if ok && r.Headers["Content-Encoding"] == "gzip" {
		bytes, size, err := compressWithGzip(r.Body)
		if err != nil {
			fmt.Println("Error compressing data: ", err.Error())
			return "", 0, err
		}
		body := string(bytes)
		return body, size, nil
	}
	body := r.Body
	return body, len(r.Body), nil

}

func compressWithGzip(str string) ([]byte, int, error) {
	data := []byte(str)
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(data)
	if err != nil {
		fmt.Println("Error compressing data:", err)
		return nil, 0, err
	}

	// Explicitly close the writer to flush all data
	if err := gzipWriter.Close(); err != nil {
		fmt.Println("Error closing gzip writer:", err)
		return nil, 0, err
	}

	compressed := buf.Bytes()
	size := len(compressed)

	// Now decompress to test
	// unzipped, _ := uncompressWithGzip(compressed)
	// fmt.Println("original: ", str, "uncompressed: ", unzipped)

	return compressed, size, nil
}

func uncompressWithGzip(data []byte) (string, error) {
	var buf bytes.Buffer
	buf.Write(data)

	gzipReader, err := gzip.NewReader(&buf)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return "", err
	}
	defer gzipReader.Close()

	uncompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		fmt.Println("Error reading uncompressed data:", err)
		return "", err
	}

	return string(uncompressedData), nil
}
