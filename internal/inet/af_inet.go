package inet

import (
	"strconv"
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
