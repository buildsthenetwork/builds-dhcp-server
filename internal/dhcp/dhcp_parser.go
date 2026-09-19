package dhcp

import "fmt"

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
DHCPREQUEST - Client message to servers either
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
	// sname 64 octets (512 bits)
	// file 128 octets (1024 bits)
	// options (variable number of bits)
}

func (msg *Message) Print() {
	fmt.Println(msg.OP)
	fmt.Println(msg.HTYPE)
	fmt.Println(msg.HLEN)
	fmt.Println(msg.HOPS)
	fmt.Println(msg.XID)
	fmt.Println(msg.SECS)
	fmt.Println(msg.FLAGS)
	fmt.Println(msg.CIADDR)
	fmt.Println(msg.YIADDR)
	fmt.Println(msg.SIADDR)
	fmt.Println(msg.GIADDR)
}
