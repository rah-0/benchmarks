package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	// Example byte arrays
	data := [][]byte{
		[]byte("Hello, World!"),
		[]byte("Newline \n inside"),
		[]byte("Another line"),
	}

	// Writing to binary file
	file, err := os.Create("/tmp/output.bin")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, bytes := range data {
		// Write length prefix (uint32)
		binary.Write(file, binary.LittleEndian, uint32(len(bytes)))
		// Write actual bytes
		file.Write(bytes)
	}

	// Reading from binary file
	file, _ = os.Open("output.bin")
	defer file.Close()

	for {
		var length uint32
		err := binary.Read(file, binary.LittleEndian, &length)
		if err != nil {
			fmt.Println("Empty")
			break // EOF
		}

		buf := make([]byte, length)
		file.Read(buf)

		fmt.Println("Decoded:", string(buf))
	}
}
