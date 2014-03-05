package ipquery

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	//"log"
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
	var length uint32
	f, err := os.Open(dataFile)
	if nil != err {
		return err
	}

	reader := bufio.NewReader(f)
	binary.Read(reader, binary.BigEndian, &length)
	data, err := ioutil.ReadAll(reader)
	if nil != err {
		return err
	}

	this.dataFile = dataFile
	this.length = length
	this.data = data
	this.cache = make(map[string]string)

	return nil
}

func (this *QueryServer) FindIp(ipstr string) (string, bool) {
	ip := net.ParseIP(ipstr)
	data := this.data
	var r []byte

	if nil == ip {
		return "", false
	}

	cachedResult, ok := this.cache[ipstr]
	if ok {
		return cachedResult, true
	}

	ip = ip.To4()

	var startPos uint32
	var endPos uint32
	var step uint32
	startPos = uint32(getStartPos(ip, data))
	endPos = this.length - 1028
	step = 8

	//log.Printf("startPos:%d,endPos:%d,step:%d", startPos, endPos, step)

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
		//log.Printf("reqIP:%d", reqIP)

		if dataIP >= reqIP {

			for j = 4; j < 7; j++ {
				resultOffsetBytes[j-4] = data[i+j]
			}
			resultOffsetBytes[3] = 0
			binary.Read(bytes.NewBuffer(resultOffsetBytes), binary.LittleEndian, &resultOffset)
			resultLength = int32(data[i+7])
			//log.Printf("FOUND!%c", IPBytes)
			//log.Printf("resultOffsetBytes:%c", resultOffsetBytes)
			//log.Printf("resultOffset:%d", resultOffset)
			//log.Printf("resultLength:%d", resultLength)

			resultOffset = resultOffset + int32(this.length) - 1024 - 4

			r = data[resultOffset : resultOffset+resultLength]
			//log.Printf("result:%s", r)

			break
		}
	}

	result := string(r)
	this.cache[ipstr] = result
	return result, true
}

func getStartPos(ip net.IP, data []byte) int32 {

	var startPos int32
	startBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		startBytes[i] = data[int(ip[0])*4+i]
	}
	binary.Read(bytes.NewBuffer(startBytes), binary.LittleEndian, &startPos)

	//log.Printf("start_bytes:%c", startBytes)
	//log.Printf("start_pos:%d", startPos)

	startPos = startPos*8 + 1024
	return startPos
}
