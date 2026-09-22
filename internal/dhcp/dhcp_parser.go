package dhcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
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
	Opcode uint8
	HType  uint8
	HLen   uint8
	Hops   uint8
	Xid    uint32
	Secs   uint16
	Flags  uint16
	CIAddr [4]byte
	YIAddr [4]byte
	SIAddr [4]byte
	GIAddr [4]byte
	// chaddr 16 octets (128 bits)
	CHAddr [16]byte
	// sname 64 octets (512 bits)
	Sname [64]byte
	// file 128 octets (1024 bits)
	File [128]byte
	// options (variable number of bits)
	MagicCookie [4]byte
	//Options     []byte
	Options []DHCPOption
}

/*
*
Prints out a message to terminal.
*/
func (m *Message) Print() {
	fmt.Println("==========PRINTING MESSAGE===========")
	fmt.Printf("OP: %s\n", parseOpcode(m.Opcode))
	fmt.Printf("HTYPE: %s\n", parseHType(m.HType))
	fmt.Printf("HLEN: %d | 0x%x\n", m.HLen, m.HLen)
	fmt.Printf("HOPS: %d | 0x%x\n", m.Hops, m.Hops)
	fmt.Printf("XID: %d | 0x%x\n", m.Xid, m.Xid)
	fmt.Printf("SECS: %d | 0x%x\n", m.Secs, m.Secs)
	fmt.Printf("FLAGS: %d | 0x%x\n", m.Flags, m.Flags)
	fmt.Printf("CIADDR: %s\n", netip.AddrFrom4(m.CIAddr))
	fmt.Printf("YIADDR: %s\n", netip.AddrFrom4(m.YIAddr))
	fmt.Printf("SIADDR: %s\n", netip.AddrFrom4(m.SIAddr))
	fmt.Printf("GIADDR: %s\n", netip.AddrFrom4(m.GIAddr))
	//fmt.Println(m.CHADDR)
	//mac := net.HardwareAddr(m.CHAddr[:m.HLen])
	fmt.Printf("CHADDR: %s\n", parseCHAddr(m.CHAddr[:m.HLen]))
	fmt.Println(m.Sname)
	fmt.Println(m.File)
	fmt.Println(m.MagicCookie)
	fmt.Println("Options: ")
	//OptionsParser(m.Options)
	printParsedOptions(m.Options)

	fmt.Println("==========END MESSAGE===========")
}

// TODO
// Need error handling for this.
// If msg length < 240, there's a problem
// If magic cookie is not correct, there's a problem
func ConstructMessage(recv []byte) Message {

	var msg Message

	msg.Opcode = recv[0]
	msg.HType = recv[1]
	msg.HLen = recv[2]
	msg.Hops = recv[3]
	msg.Xid = binary.BigEndian.Uint32(recv[4:8])
	msg.Secs = binary.BigEndian.Uint16(recv[8:10])
	msg.Flags = binary.BigEndian.Uint16(recv[10:12])
	// msg.CIADDR = binary.BigEndian.Uint32(recv[12:16])
	copy(msg.CIAddr[:], recv[12:16])
	copy(msg.YIAddr[:], recv[16:20])
	copy(msg.SIAddr[:], recv[20:24])
	copy(msg.GIAddr[:], recv[24:28])

	copy(msg.CHAddr[:], recv[28:44]) // TODO copy works, but probably not what we want
	copy(msg.Sname[:], recv[44:108])
	copy(msg.File[:], recv[108:236])
	copy(msg.MagicCookie[:], recv[236:240])
	//msg.Options = recv[240:] // This might be a copy by instance, not value
	msg.Options = OptionsParser(recv[240:])

	return msg
}

func parseOpcode(op uint8) string { // eventually string, err
	switch op {
	case 1:
		return "BOOTREQUEST (1)"
	case 2:
		return "BOOTREPLY (2)"
	default:
		return "UNKNOWN"
	}
}

func parseHType(htype uint8) string {
	switch htype {
	case 1:
		return "Ethernet (1)"
	case 2:
		return "Experimental Ethernet (2)"
	case 3:
		return "Amateur Radio AX.25 (3)"
	case 4:
		return "Proteon ProNET Token Ring (4)"
	case 5:
		return "Chaos (5)"
	case 6:
		return "IEEE 802 Networks (6)"
	case 7:
		return "ARCNET (7)"
	case 8:
		return "Hyperchannel (8)"
	case 9:
		return "Lanstar (9)"
	default:
		return "UNKNOWN - (" + string(htype) + ")"
	}
}

func parseCHAddr(haddr []byte) string {
	mac := net.HardwareAddr(haddr)

	return mac.String()

}
