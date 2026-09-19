package main

import (
	"fmt"

	"github.com/buildsthenetwork/builds-dhcp-server/internal/dhcp"
)

func main() {
	fmt.Println("Hello DHCP!")

	dhcp.BindServer()

}
