package ipquery

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	//"log"
	"os"
)

type Parser_17mon struct {
}

func (this *Parser_17mon) ParseData(dataFile string) ([]DataItem, error) {

	/* --Read the file */
	f, err := os.Open(dataFile)
	if nil != err {
		return nil, err
	}

	reader := bufio.NewReader(f)
	data, err := ioutil.ReadAll(reader)
	if nil != err {
		return nil, err
	}

	/* --Get record length */
	var length uint32
	lengthBytes := data[0:4]
	binary.Read(bytes.NewBuffer(lengthBytes), binary.BigEndian, &length)

	/* --Calc start position and end position */
	startPos, readBytesCnt := binary.Varint(data[4:8])
	if readBytesCnt != 1 {
		//log.Fatalf("[ParseData] Get start position failed,readBytesCnt:%d ", readBytesCnt)
		return nil, errors.New("[ParseData] Get start position failed")
	}
	endPos := int64(length - 1024)
	//log.Printf("[ParseData] StartPos:%d,EndPos:%d", startPos, endPos)

	/* --Parse items */
	items := []DataItem{}
	for i := startPos + 1028; i < endPos; i += 8 {
		var ip []byte
		var offset uint32
		var strlength uint32
		var result string

		/* --- Ip */
		ip = data[i+0 : i+4]
		/* --- Offset */
		bufOffset := make([]byte, 4)
		bufOffset[0] = data[i+4]
		bufOffset[1] = data[i+5]
		bufOffset[2] = data[i+6]
		bufOffset[3] = 0
		binary.Read(bytes.NewBuffer(bufOffset), binary.LittleEndian, &offset)
		offset = offset + length - 1024
		/* --- StrLength */
		strlength = uint32(data[i+7])
		/* --- Get result string */
		result = string(data[offset : offset+strlength])
		/* --- Append to items */
		items = append(items, makeItemData(ip, result))

		//log.Printf("[ParseData] Got Item: {ip:%c,offset:%d,strlength:%s}", ip, offset, strlength)
		//log.Printf("[ParseData] Addr: %s", result)
	}

	//log.Printf("[ParseData]item data length:%d", len(this.items))

	return items, nil

}

func makeItemData(ip []byte, address string) DataItem {
	address_struct := getAddressStructByAddressString(address)
	item := DataItem{Ip: ip2long(ip), Address: address_struct}
	return item
}
