package inet

import (
	"strconv"
	"strings"
)

func Int_to_addr(iaddr uint32) string {

	octet4 := strconv.Itoa(int(iaddr & 0xff))
	iaddr = iaddr >> 8
	octet3 := strconv.Itoa(int(iaddr & 0xff))
	iaddr = iaddr >> 8
	octet2 := strconv.Itoa(int(iaddr & 0xff))
	iaddr = iaddr >> 8
	octet1 := strconv.Itoa(int(iaddr & 0xff))

	return octet1 + "." + octet2 + "." + octet3 + "." + octet4
}

func Addr_to_int(addr string) uint32 {
	spl := strings.Split(addr, ".")

	// fmt.Println(spl)
	octet1, err := strconv.Atoi(spl[0])
	if err != nil {
		panic(err)
	}
	octet2, err := strconv.Atoi(spl[1])
	if err != nil {
		panic(err)
	}
	octet3, err := strconv.Atoi(spl[2])
	if err != nil {
		panic(err)
	}
	octet4, err := strconv.Atoi(spl[3])
	if err != nil {
		panic(err)
	}

	octet1 = octet1 << 24
	octet2 = octet2 << 16
	octet3 = octet3 << 8

	iaddr := octet1 + octet2 + octet3 + octet4

	return uint32(iaddr)

}
