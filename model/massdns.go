package model

type DNSQuery struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Class  string `json:"class"`
	Status string `json:"status"`
	RxTs   int64  `json:"rx_ts"`
	Data   struct {
		DataAuthorities []DataAuthorities `json:"authorities,omitempty"`
		DataAnswers     []DataAnswers     `json:"answers,omitempty"`
	} `json:"data,omitempty"`
	Flags    []string `json:"flags"`
	Resolver string   `json:"resolver"`
	Proto    string   `json:"proto"`
}
type DataAuthorities struct {
	TTL   int    `json:"ttl"`
	Type  string `json:"type"`
	Class string `json:"class"`
	Name  string `json:"name"`
	Data  string `json:"data"`
}
type DataAnswers struct {
	TTL   int    `json:"ttl"`
	Type  string `json:"type"`
	Class string `json:"class"`
	Name  string `json:"name"`
	Data  string `json:"data"`
}
