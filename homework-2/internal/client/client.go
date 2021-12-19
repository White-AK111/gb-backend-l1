package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// InitClient default client from training manual
func InitClient() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		_, err = io.Copy(os.Stdout, conn)
		if err != nil {
			return
		}
	}()
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		return
	}
	fmt.Printf("%s: exit", conn.LocalAddr())
}
