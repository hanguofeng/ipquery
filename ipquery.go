package ipquery

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	//"log"
	"net"
	"os"
	"strings"
)

type QueryServer struct {
	dataFile string
	length   uint32
	data     []byte
	items    []DataItem
}

type DataItem struct {
	Ip      uint32
	Address *Address
}

type Address struct {
	Country  string
	Province string
	City     string
}

func (this *QueryServer) Load(dataFile string) error {
	f, err := os.Open(dataFile)
	if nil != err {
		return err
	}

	reader := bufio.NewReader(f)
	data, err := ioutil.ReadAll(reader)
	if nil != err {
		return err
	}

	this.dataFile = dataFile
	this.data = data

	this.length = this.getLength()
	this.ParseData()

	return nil
}

func (this *QueryServer) ParseData() error {
	data := this.data
	startPos, readBytesCnt := binary.Varint(data[4:8])
	if readBytesCnt != 1 {
		//log.Fatalf("[ParseData] Get start position failed,readBytesCnt:%d ", readBytesCnt)
		return errors.New("[ParseData] Get start position failed")
	}
	endPos := int64(this.length - 1024)
	//log.Printf("[ParseData] StartPos:%d,EndPos:%d", startPos, endPos)

	this.initItemData()
	for i := startPos + 1028; i < endPos; i += 8 {
		var ip []byte
		var offset uint32
		var length uint32
		var result string
		bufOffset := make([]byte, 4)

		ip = data[i+0 : i+4]
		bufOffset[0] = data[i+4]
		bufOffset[1] = data[i+5]
		bufOffset[2] = data[i+6]
		bufOffset[3] = 0
		binary.Read(bytes.NewBuffer(bufOffset), binary.LittleEndian, &offset)
		offset = offset + this.length - 1024
		length = uint32(data[i+7])
		result = string(data[offset : offset+length])
		this.addItemData(ip, result)

		//log.Printf("[ParseData] Got Item: {ip:%c,offset:%d,length:%s}", ip, offset, length)
		//log.Printf("[ParseData] Addr: %s", result)
	}

	//log.Printf("[ParseData]item data length:%d", len(this.items))
	return nil

}

func (this *QueryServer) FindIp(ipstr string) (*Address, bool) {
	ip := net.ParseIP(ipstr)
	if nil == ip {
		return nil, false
	}
	ip = ip.To4()
	iplong := ip2long(ip)

	itemsCount := len(this.items)

	low := 0
	high := itemsCount - 1

	if iplong == this.items[low].Ip {
		return this.items[low].Address, true
	}

	if iplong == this.items[high].Ip {
		return this.items[high].Address, true
	}

	for low <= high {
		mid := low + ((high - low) / 2)

		//log.Printf("[FindIp2]low:%d,high:%d,mid:%d,midip:%d,findip:%d", low, high, mid, this.items[mid].Ip, iplong)

		if this.items[mid].Ip >= iplong {
			if mid >= 1 {
				if this.items[mid-1].Ip <= iplong {
					return this.items[mid].Address, true
				}
			} else {
				return this.items[mid].Address, true
			}
		}

		if this.items[mid].Ip > iplong {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return nil, false
}

func (this *QueryServer) getLength() uint32 {
	var length uint32
	data := this.data
	lengthBytes := data[0:4]
	binary.Read(bytes.NewBuffer(lengthBytes), binary.BigEndian, &length)
	//log.Printf("[getLength]length:%d", length)
	return length
}

func (this *QueryServer) initItemData() {
	this.items = []DataItem{}
}

func (this *QueryServer) addItemData(ip []byte, address string) {
	address_struct := getAddressStructByAddressString(address)
	item := DataItem{Ip: ip2long(ip), Address: address_struct}
	this.items = append(this.items, item)
}

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
