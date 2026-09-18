package main

import (
	"fmt"

	"github.com/buildsthenetwork/builds-dhcp-server/internal/inet"
)

func main() {
	fmt.Println("Hello World! DHCP")
	addr := inet.Int_to_addr(0b11000000_10101000_00001010_00000001)
	fmt.Println(addr)
}
