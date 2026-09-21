package dhcp

import (
	"fmt"
)

/*

Options requirements are being built according to RFC 2132.

https://www.rfc-editor.org/info/rfc2132/

*/

// A structure to map Option codes to Option Names
var optionNames = map[int]string{
	0:  "Pad Option",
	1:  "Subnet Mask",
	2:  "Time Offset",
	3:  "Router",
	4:  "Time Server",
	5:  "Name Server",
	6:  "Domain Name Server",
	7:  "Log Server",
	8:  "Cookie Server",
	9:  "LPR Server", // Line Printer Servers
	10: "Impress Server",
	11: "Resource Location Server",
	12: "Host Name",
	13: "Boot File Size",
	14: "Merit Dump File",
	15: "Domain Name",
	16: "Swap Server",
	17: "Root Path",
	18: "Extensions Path",
	19: "IP Forwarding Enable/Disable",
	20: "Non-Local Source Routing Enable/Disable",
	21: "Policy Filter",
	22: "Maximum Datagram Reassembly Size",
	23: "Default IP Time-to-live",
	24: "Path MTU Aging Timeout",
	25: "Path MTU Plateau Table",
	26: "Interface MTU",
	27: "All Subnets are Local",
	28: "Broadcast Address",
	29: "Perform Mask Discovery",
	30: "Mask Supplier",
	31: "Perform Router Discovery",
	32: "Router Solicitation Address",
	33: "Static Route",
	34: "Trailer Encapsulation",
	35: "ARP Cache Timeout",
	36: "Ethernet Encapsulation",
	37: "TCP Default TTL",
	38: "TCP Keepalive Internal",
	39: "TCP Keepalive Garbage",
	40: "Network Information Service Domain",
	41: "Network Information Servers",
	42: "Network Time Protocol Servers",
	43: "Vendor Specific Information",
	44: "NetBIOS over TCP/IP Name Server",
	45: "NetBIOS over TCP/IP Datagram Distribution Server",
	46: "NetBIOS over TCP/IP Node Type",
	47: "NetBIOS over TCP/IP Scope",
	48: "X Window System Font Server",
	49: "X Window System Display Manager",
	50: "Requested IP Address",
	51: "IP Address Lease Time",
	52: "Option Overload",
	53: "DHCP Message Type",
	54: "Server Identifier",
	55: "Parameter Request List",
	56: "Message",
	57: "Maximum DHCP Message Size",
	58: "Renewal (T1) Time Value",
	59: "Rebinding (T2) Time Value",
	60: "Vendor Class Identifier",
	61: "Client-Identifier",

	64: "Network Information Service+ Domain",
	65: "Network Information Service+ Servers",

	68: "Mobile IP Home Agent",
	69: "Simple Mail Transport Protocol (SMTP) Server",
	70: "Post Office Protocol (POP3) Server",
	71: "Network News Transport Protocol (NNTP) Server",
	72: "Default World Wide Web (WWW) Server",
	73: "Default Finger Server",
	74: "Default Internet Relay Chat (IRC) Server",
	75: "StreetTalk Server",
	76: "StreetTalk Directory Assistance (STDA) Server",

	255: "End Option",
}

type DHCPOption struct {
	code   uint8
	length uint8
	data   []byte
}

func (option *DHCPOption) Print() {
	fmt.Println("    ------------")
	fmt.Printf("Option: %s (%d)\n", optionNames[int(option.code)], option.code)
	fmt.Printf("Length: %d\n", option.length)
	fmt.Print("Message Data: ")
	fmt.Print(option.data)
	fmt.Print("\n")
	parseOptionData(option)
	fmt.Println("    ------------")
}

func OptionsParser(options []byte) {
	// options 0 and 255 are fixed length. I think 255 is typically the end...
	// 128 to 254 are reserverd for site-specific options. However, I do know option 150 is Cisco TFTP

	// 0 is the pad option
	// 255 is the end option, all subsequent options are pad options

	// Subnet Mask option
	// If both the subnet mask and the router option are specified in a reply, the subnet mask option MUST be first.
	dhcp_options := []DHCPOption{}
	for i := 0; i < len(options); i++ {
		// append
		// slice_name = append(slice_name, element1, element2, ...)
		code := options[i]

		if code == 0 {
			continue // skip the iteration
		}

		if code == 255 {
			fmt.Println("Option 255. End Parsing!")
			break
		}

		// TODO - include validation for malformed packets
		if i+1 > len(options) {
			break
		}

		length := int(options[i+1])
		start := i + 2        // inclusive
		end := start + length // exclusive. At end of loop, i will equal the end

		// the options[start:end] slice can potentially get messy. Look for a different way.
		option := DHCPOption{code: code, length: uint8(length), data: options[start:end]}
		dhcp_options = append(dhcp_options, option)
		i = end - 1
	}

	fmt.Println("-------- Options Parser --------")
	//option := DHCPOption{code: options[0], length: options[1], data: options[2:3]}
	//option.Print()
	for i := 0; i < len(dhcp_options); i++ {
		dhcp_options[i].Print()
	}
	fmt.Println("--------   END   Parser --------")
}

func parseOptionData(option *DHCPOption) {

	switch option.code {
	case 53:
		option53(option.data)
	default:
		fmt.Printf("--Option %d not implemented.\n", option.code)
	}
}

// 53 - DHCP Message Type
// Length is always 1
func option53(data []byte) {
	msg_type := data[0]
	msg_str := "Message Type: "
	switch msg_type {
	case 1:
		msg_str += "DHCPDISCOVER"
	case 2:
		msg_str += "DHCPOFFER"
	case 3:
		msg_str += "DHCPREQUEST"
	case 4:
		msg_str += "DHCPDECLINE"
	case 5:
		msg_str += "DHCPACK"
	case 6:
		msg_str += "DHCPNAK"
	case 7:
		msg_str += "DHCPRELEASE"
	case 8:
		msg_str += "DHCPINFORM"
	default:
		msg_str += "UNKNOWN DHCP MESSAGE TYPE"
	}
	fmt.Println(msg_str)
}
