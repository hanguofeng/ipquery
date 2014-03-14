package ipquery

type DataItem struct {
	Ip      uint32
	Address *Address
}

type Address struct {
	Country  string
	Province string
	City     string
}
