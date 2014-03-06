package ipquery

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	//	"log"
	"net"
	"os"
)

type QueryServer struct {
	dataFile string
	length   uint32
	data     []byte
	cache    map[string]string
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
	this.cache = make(map[string]string)

	return nil
}

func (this *QueryServer) FindIp(ipstr string) (string, bool) {
	ip := net.ParseIP(ipstr)
	if nil == ip {
		return "", false
	}
	ip = ip.To4()

	cachedResult, ok := this.cache[ipstr]
	if ok {
		return cachedResult, true
	}

	startPos := uint32(this.getStartPos(ip))
	endPos := uint32(this.length - 1028)
	resultStr := this.findResultString(ip, startPos, endPos)
	this.cache[ipstr] = resultStr
	return resultStr, true
}

func (this *QueryServer) getStartPos(ip net.IP) int32 {
	var startPos int32
	data := this.data
	startBytes := data[int(ip[0])*4 : int(ip[0])*4+4]
	binary.Read(bytes.NewBuffer(startBytes), binary.LittleEndian, &startPos)
	//	log.Printf("[getStartPos]start_bytes:%c,start_pos:%d", startBytes, startPos)
	startPos = startPos*8 + 1024 - 4 //4 means the first 4 bytes are info about length
	return startPos
}

func (this *QueryServer) getLength() uint32 {
	var length uint32
	data := this.data
	lengthBytes := data[0:4]
	binary.Read(bytes.NewBuffer(lengthBytes), binary.BigEndian, &length)
	//log.Printf("[getLength]length:%d", length)
	return length
}

func (this *QueryServer) findResultString(ip net.IP, startPos uint32, endPos uint32) string {

	var r []byte
	data := this.data

	step := uint32(8)

	for i := startPos; i < endPos; i += step {
		var dataIP int32
		var reqIP int32
		var j uint32
		var resultOffset int32
		var resultLength int32

		IPBytes := make([]byte, 4)
		resultOffsetBytes := make([]byte, 4)
		for j = 0; j < 4; j++ {
			IPBytes[j] = data[i+j]
		}

		binary.Read(bytes.NewBuffer(IPBytes), binary.BigEndian, &dataIP)
		binary.Read(bytes.NewBuffer(ip), binary.BigEndian, &reqIP)

		//log.Printf("startPos:%d", startPos)
		//log.Printf("IPBytes:%c", IPBytes)
		//log.Printf("dataIP:%d", dataIP)
		//log.Printf("reqIP:%d", reqIP)0xc200088000

		if dataIP >= reqIP {

			for j = 4; j < 7; j++ {
				resultOffsetBytes[j-4] = data[i+j]
			}

			binary.Read(bytes.NewBuffer(resultOffsetBytes), binary.LittleEndian, &resultOffset)
			resultLength = int32(data[i+7])
			//log.Printf("FOUND!%c", IPBytes)
			//log.Printf("resultOffsetBytes:%c", resultOffsetBytes)
			//log.Printf("resultOffset:%d", resultOffset)
			//log.Printf("resultLength:%d", resultLength)

			resultOffset = resultOffset + int32(this.length) - 1024

			r = data[resultOffset : resultOffset+resultLength]
			//log.Printf("result:%s", r)

			break
		}
	}

	result := string(r)
	return result
}
