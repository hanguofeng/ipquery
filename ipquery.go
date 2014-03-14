package ipquery

import (
	//"log"
	"net"
)

type QueryServer struct {
	dataFile string
	length   uint32
	data     []byte
	items    []DataItem
}

func (this *QueryServer) Load(dataFile string) error {
	parser := &Parser_17mon{}
	items, err := parser.ParseData(dataFile)
	if nil != err {
		return err
	} else {
		this.items = items
		return nil
	}

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
