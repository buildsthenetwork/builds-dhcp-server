package dhcp

import (
	"fmt"
)

// every option has a name and pre-defined length.

var optionNames = map[int]string{
	0:   "Pad Option",
	1:   "Subnet Mask",
	12:  "Host Name",
	50:  "Requested IP Address",
	51:  "IP Address Lease Time",
	53:  "DHCP Message Type",
	55:  "Parameter Request List",
	57:  "Maximum DHCP Size",
	61:  "Client Identifier",
	255: "End Option",
}

type DHCPOption struct {
	code uint8
	//name   string
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
	// my first instinct is going to be: [(tag octet) (len octet) (number of items given the len)]
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
		//fmt.Println("DHCPDISCOVER")
		msg_str += "DHCPDISCOVER"
	case 2:
		//fmt.Println("DHCPOFFER")
		msg_str += "DHCPOFFER"
	case 3:
		//fmt.Println("DHCPREQUEST")
		msg_str += "DHCPREQUEST"
	case 4:
		//fmt.Println("DHCPDECLINE")
		msg_str += "DHCPDECLINE"
	case 5:
		//fmt.Println("DHCPACK")
		msg_str += "DHCPACK"
	case 6:
		//fmt.Println("DHCPNAK")
		msg_str += "DHCPNAK"
	case 7:
		//fmt.Println("DHCPRELEASE")
		msg_str += "DHCPRELEASE"
	case 8:
		//fmt.Println("DHCPINFORM")
		msg_str += "DHCPINFORM"
	default:
		//fmt.Println("UNKNOWN DHCP MESSAGE TYPE")
		msg_str += "UNKNOWN DHCP MESSAGE TYPE"
	}
	fmt.Println(msg_str)
}
