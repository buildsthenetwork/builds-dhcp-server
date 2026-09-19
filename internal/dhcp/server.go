package dhcp

import (
	"fmt"
	"net"
)

func BindServer() {
	addr, err := net.ResolveUDPAddr("udp", ":67")

	if err != nil {
		fmt.Println("Failed to resolve address: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to listen on UDP:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server is listening on port 67...")

	buffer := make([]byte, 1024) // I'm not positive on the size we want
	for {
		// Read data sent by client
		n, clientAddr, err := conn.ReadFromUDP(buffer) // is this blocking?

		if err != nil {
			fmt.Println("Error reading data:", err)
			continue
		}

		message := buffer[:n] // is n the length?
		fmt.Printf("Recieved message from %s...\n", clientAddr)
		fmt.Println(message)
	}

}
