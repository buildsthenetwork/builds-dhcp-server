package dhcp_test

import (
	"testing"
)

func TestParseDHCPDiscover(t *testing.T) {
	packet := make([]byte, 240)

	// BOOTP header
	packet[0] = 1 // BOOTREQUEST
	packet[1] = 1 // HTYPE ETHERNET
	packet[2] = 6 // HLEN 6, likely a MAC Address
	packet[3] = 0 // HOPS 0

	// Transaction ID
	copy(packet[4:8], []byte{0x12, 0x34, 0x56, 0x78})

	// Flags
	packet[10] = 0x80 // Ah, smart. Instead of 10:11, you only needed to manipulate 10.

	// Client hardware address
	copy(packet[28:34], []byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55,
	})

	// DHCP magic cookie
	copy(packet[236:240], []byte{
		0x63, 0x82, 0x53, 0x63,
	})

	// DHCP options
	options := []byte{
		53, 1, 1,
		61, 7, 1, 0, 17, 34, 51, 68, 85,
		55, 4, 1, 3, 6, 15,
		12, 4, 84, 69, 83, 84,
		200, 3, 0xDE, 0xAD, 0xBE,
		255,
	}

	packet = append(packet, options...)

	// Pad to 300 bytes
	for len(packet) < 300 {
		packet = append(packet, 0)
	}

	//fmt.Printf("Packet length: %d bytes\n", len(packet))
	//fmt.Println(hex.Dump(packet))

	// Feed packet into your parser here.
}
