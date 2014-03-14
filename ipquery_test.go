package ipquery

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const DATA_FILE = "./data/17monipdb.dat"

func TestLoad(t *testing.T) {
	s := new(QueryServer)
	var err error
	err = s.Load(DATA_FILE)

	if nil != err {
		t.Errorf("Load Error:%s", err.Error())
	}

	err = s.Load("NOT_EXIST_FILE")

	if nil == err {
		t.Errorf("Load Error should be error because of not exist file")
	}
}

func TestFindIp(t *testing.T) {
	s := new(QueryServer)
	err := s.Load(DATA_FILE)

	if nil != err {
		t.Errorf("Load Error:%s", err.Error())
	}

	var ipstr string
	var address *Address
	var ok bool

	ipstr = getRandomIpStr()
	address, ok = s.FindIp(ipstr)
	if !ok {
		t.Errorf("FindIp Failed")
	}
	t.Logf("ipstr:%s,address:%s", ipstr, address)

	ipstr = "127.0.0.1"
	address, ok = s.FindIp(ipstr)
	if !ok {
		t.Errorf("FindIp Failed")
	}
	t.Logf("ipstr:%s,address:%s", ipstr, address)

	ipstr = "0.0.0.1"
	address, ok = s.FindIp(ipstr)
	if !ok {
		t.Errorf("FindIp Failed")
	}
	t.Logf("ipstr:%s,address:%s", ipstr, address)

	ipstr = "202.194.34.229"
	address, ok = s.FindIp(ipstr)
	if !ok {
		t.Errorf("FindIp Failed")
	}
	t.Logf("ipstr:%s,address:%s", ipstr, address)

}

func BenchmarkLoad(b *testing.B) {
	s := new(QueryServer)

	for i := 0; i < b.N; i++ {
		err := s.Load(DATA_FILE)
		if nil != err {
			b.Errorf("Load Error:%s", err.Error())
		}
	}
}

func BenchmarkFindIp(b *testing.B) {
	s := new(QueryServer)
	err := s.Load(DATA_FILE)
	if nil != err {
		b.Errorf("Load Error:%s", err.Error())
	}

	ipstr := getRandomIpStr()
	for i := 0; i < b.N; i++ {
		_, ok := s.FindIp(ipstr)
		if !ok {
			b.Errorf(DATA_FILE)
		}
	}
}

func BenchmarkFindIpRandom(b *testing.B) {
	s := new(QueryServer)
	err := s.Load(DATA_FILE)
	if nil != err {
		b.Errorf("Load Error:%s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ipstr := getRandomIpStr()
		_, ok := s.FindIp(ipstr)
		if !ok {
			b.Errorf("FindIp Failed")
		}
	}
}

func getRandomIpStr() string {
	return fmt.Sprintf("%d.%d.%d.%d", rnd(0, 254), rnd(0, 254), rnd(0, 254), rnd(0, 254))
}

func rnd(from, to int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(to+1-from) + from
}
