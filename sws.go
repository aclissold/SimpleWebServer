package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// Listen for TCP requests on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Start an infinite loop to accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a different goroutine to avoid blocking
		go handleConnection(conn)
	}
}

// handleConnection naïvely assumes the connection is sending an HTML GET request
// for a file and attempts to serve up that file.
func handleConnection(conn net.Conn) {
	// Close the connection when we're done with it and give it a 5-second time limit
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read the first line of the request and parse the filename from it
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	requestLine := strings.Split(scanner.Text(), " ")
	if len(requestLine) < 2 {
		log.Println("invalid request")
		return
	}
	filename := "." + requestLine[1] // prepend . to avoid root!

	// Open, read, and write the file to the connection
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		_, pathError := err.(*os.PathError)
		if pathError {
			conn.Write(append(notFoundHeader, []byte("404 Not Found")...))
			return
		}
		log.Println(err)
		return
	}
	buf := make([]byte, 1024)
	file.Read(buf)
	conn.Write(append(okHeader, buf...))
}

var okHeader = []byte("HTTP/1.1 200 OK\r\nServer : Go\r\n\r\n")
var notFoundHeader = []byte("HTTP/1.1 404 Not Found\r\nServer : Go\r\n\r\n")
