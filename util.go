package ipquery

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func getAddressStructByAddressString(address string) *Address {
	address_arr := strings.Split(address, "\t")
	country := ""
	province := ""
	city := ""

	if len(address_arr) > 0 {
		country = address_arr[0]
	}
	if len(address_arr) > 1 {
		province = address_arr[1]
	}
	if len(address_arr) > 2 {
		city = address_arr[2]
	}

	address_struct := &Address{Country: country, Province: province, City: city}
	return address_struct
}

func ip2long(ip []byte) uint32 {
	var iplong uint32
	binary.Read(bytes.NewBuffer(ip), binary.BigEndian, &iplong)
	return iplong
}
