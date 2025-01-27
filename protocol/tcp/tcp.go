package main

import (
	"bufio"
	"log"
	"net"
)

var (
	address = "localhost:55555"
	server  net.Listener
	client  *TCPClient
)

func tcpStartServer() (err error) {
	server, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("TCP server listening on %s", address)

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}
			go handleTCPConnection(conn)
		}
	}()

	return nil
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("Connection closed by client")
			} else {
				log.Printf("Connection error: %v", err)
			}
			return
		}

		// Process the received bytes (optional)
		_ = buffer[:n] // Do nothing with the data
	}
}

type TCPClient struct {
	conn   net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
}

func (c *TCPClient) connect() error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.writer = bufio.NewWriter(conn)
	c.reader = bufio.NewReader(conn)

	return nil
}

func (c *TCPClient) send(data []byte) error {
	// Write the raw bytes
	_, err := c.writer.Write(data)
	if err != nil {
		return err
	}

	// Flush the writer to ensure data is sent immediately
	return c.writer.Flush()
}

func (c *TCPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
