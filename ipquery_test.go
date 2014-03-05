package ipquery

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestLoad(t *testing.T) {
	s := new(QueryServer)
	err := s.Load("./data/17monipdb.dat")

	if nil != err {
		t.Errorf("Load Error:%s", err.Error())
	}

	ipstr := getRandomIpStr()
	address, ok := s.FindIp(ipstr)
	if !ok {
		t.Errorf("FindIp Failed")
	}
	t.Logf("ipstr:%s,address:%s", ipstr, address)
}

func BenchmarkLoad(b *testing.B) {
	s := new(QueryServer)

	for i := 0; i < b.N; i++ {
		err := s.Load("./data/17monipdb.dat")
		if nil != err {
			b.Errorf("Load Error:%s", err.Error())
		}
	}
}

func BenchmarkFindIp(b *testing.B) {
	s := new(QueryServer)
	err := s.Load("./data/17monipdb.dat")
	if nil != err {
		b.Errorf("Load Error:%s", err.Error())
	}

	ipstr := getRandomIpStr()
	for i := 0; i < b.N; i++ {
		address, ok := s.FindIp(ipstr)
		if !ok {
			b.Errorf("FindIp Failed")
		}
		b.Logf("ipstr:%s,address:%s", ipstr, address)
	}
}

func BenchmarkFindIpRandom(b *testing.B) {
	s := new(QueryServer)
	err := s.Load("./data/17monipdb.dat")
	if nil != err {
		b.Errorf("Load Error:%s", err.Error())
	}

	for i := 0; i < b.N; i++ {
		ipstr := getRandomIpStr()
		address, ok := s.FindIp(ipstr)
		if !ok {
			b.Errorf("FindIp Failed")
		}
		b.Logf("ipstr:%s,address:%s", ipstr, address)
	}
}

func getRandomIpStr() string {
	return fmt.Sprintf("%d.%d.%d.%d", rnd(0, 254), rnd(0, 254), rnd(0, 254), rnd(0, 254))
}

func rnd(from, to int) int {
	return rand.Intn(to+1-from) + from
}
