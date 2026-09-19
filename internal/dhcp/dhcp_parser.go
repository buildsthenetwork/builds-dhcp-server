package dhcp

import (
	"encoding/binary"
	"fmt"

	"github.com/buildsthenetwork/builds-dhcp-server/internal/inet"
)

/*

op (1) - message op code / message type
htype (1) - hardware address type (MAC address?)
hlen (1) - hardware address length '6' for 10mb ethernet
hops (1) - client sets to 0, optionally used by relay agents when booting via a relay agent
xid (4) - transaction ID, randomly chosen by client
secs (2) - Filled in by client, seconds elapsed since client began address aquisition or renewal process
flags (2)
ciaddr (4) - Client IP address; only filled in if client is in BOUND, RENEW, or REBINDING state and can respond to ARP requests
yiaddr (4) - 'your' (client) IP address (given by server in DHCPOFFER)
				Server need not reserve this, but probably should.
				Server may probe the offered address with an ICMP Echo Request
siaddr (4) - IP address of next server to use in bootstrap;
			returned in DHCPOFFER, DHCPACK by server.
giaddr (4) - Relay agent IP address, used in booting via a relay agent
chaddr (16) - client hardware address
sname (64) - Optional server hostname, null terminated string.
file (128) - Boot file name, null terminated string; "generic"
			name or null in DHCPDISCOVER, fully qualified directory-path name in DHCPOFFER.
options (var) Optional parameter field.


options to look for DHCP message type
The first four octets of the 'options' field of the DHCP message contain
the (decimal) values of 99, 130, 83 annd 99, respectively (this is the same
magic cookie as defined in RFC 1497)

DHCP message types are Option 53
DHCPDISCOVER (1) - Client broadcast to locate available servers.
DHCPOFFER - Server to client in response to DHCPDISCOVER with offer of config parameters.
DHCPREQUEST (3) - Client message to servers either
			(a) requesting offered parameters from one server, and implicitly declining offers from all others
			(b) confirming correctness of previously allocated address affter e.g. system reboot
			(c) extending the lease on a particular network address
DHCPACK - Server to client with config parameters, including committed network address.
DHCPNAK - Server to client indicating client's notion of network address is incorrect (e.g. client has moved to new subnet) or client's lease as expired
DHCPDECLINE - Client to server indicating network address is already in use.
DHCPRELEASE - Client to server relinquishing network address and cancelling remaining lease.
DHCPINFORM - Client to server, asking only for local config paramters. Client has externally configured network address.

*/

type Message struct {
	OP     uint8
	HTYPE  uint8
	HLEN   uint8
	HOPS   uint8
	XID    uint32
	SECS   uint16
	FLAGS  uint16
	CIADDR uint32
	YIADDR uint32
	SIADDR uint32
	GIADDR uint32
	// chaddr 16 octets (128 bits)
	CHADDR [16]byte
	// sname 64 octets (512 bits)
	SNAME [64]byte
	// file 128 octets (1024 bits)
	FILE [128]byte
	// options (variable number of bits)
	MAGIC_COOKIE [4]byte
	OPTIONS      []byte
}

func (m *Message) Print() {
	fmt.Println("==========PRINTING MESSAGE===========")

	fmt.Printf("OP: %d 0x%x\n", m.OP, m.OP)
	fmt.Printf("HTYPE: %d 0x%x\n", m.HTYPE, m.HTYPE)
	fmt.Printf("HLEN: %d 0x%x\n", m.HLEN, m.HLEN)
	fmt.Printf("HOPS: %d 0x%x\n", m.HOPS, m.HOPS)
	fmt.Printf("XID: %d 0x%x\n", m.XID, m.XID)
	fmt.Printf("SECS: %d 0x%x\n", m.SECS, m.SECS)
	fmt.Printf("FLAGS: %d 0x%x\n", m.FLAGS, m.FLAGS)
	fmt.Printf("CIADDR: %s\n", inet.Int_to_addr(m.CIADDR))
	fmt.Printf("YIADDR: %s\n", inet.Int_to_addr(m.YIADDR))
	fmt.Printf("SIADDR: %s\n", inet.Int_to_addr(m.SIADDR))
	fmt.Printf("GIADDR: %s\n", inet.Int_to_addr(m.GIADDR))
	fmt.Println(m.CHADDR)
	fmt.Println(m.SNAME)
	fmt.Println(m.FILE)
	fmt.Println(m.MAGIC_COOKIE)
	fmt.Println(m.OPTIONS)
	fmt.Println("==========END MESSAGE===========")
}

// TODO
// Need error handling for this.
// If msg length < 240, there's a problem
// If magic cookie is not correct, there's a problem
func ParseMessage(recv []byte) Message {

	var msg Message

	msg.OP = recv[0]
	msg.HTYPE = recv[1]
	msg.HLEN = recv[2]
	msg.HOPS = recv[3]
	msg.XID = binary.BigEndian.Uint32(recv[4:8])
	msg.SECS = binary.BigEndian.Uint16(recv[8:10])
	msg.FLAGS = binary.BigEndian.Uint16(recv[10:12])
	msg.CIADDR = binary.BigEndian.Uint32(recv[12:16])
	msg.YIADDR = binary.BigEndian.Uint32(recv[16:20])
	msg.SIADDR = binary.BigEndian.Uint32(recv[20:24])
	msg.GIADDR = binary.BigEndian.Uint32(recv[24:28])
	copy(msg.CHADDR[:], recv[28:44]) // TODO copy works, but probably not what we want
	copy(msg.SNAME[:], recv[44:108])
	copy(msg.FILE[:], recv[108:236])
	copy(msg.MAGIC_COOKIE[:], recv[236:240])
	msg.OPTIONS = recv[240:] // This might be a copy by instance, not value

	return msg
}
