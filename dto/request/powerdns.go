package request

type ZoneCreate struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Nameservers []string `json:"nameservers"`
}

type ZoneDelete struct {
	ZoneId string `json:"zone_id"`
}
type RecordGet struct {
	ZoneId string `json:"zone_id"`
}
type RecordData struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
	Name     string `json:"name"`
	TTL      int    `json:"ttl"`
	Type     string `json:"type"`
}
type RecordModify struct {
	Rrsets []struct {
		Comments   []interface{} `json:"comments"`
		Name       string        `json:"name"`
		Type       string        `json:"type"`
		TTL        int           `json:"ttl"`
		Changetype string        `json:"changetype"`
		Records    []struct {
			Content  string `json:"content"`
			Disabled bool   `json:"disabled"`
			Name     string `json:"name"`
			TTL      int    `json:"ttl"`
			Type     string `json:"type"`
		} `json:"records"`
	} `json:"rrsets"`
}

type RecordDelete struct {
	Rrsets []struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		TTL        int    `json:"ttl"`
		Changetype string `json:"changetype"`
	} `json:"rrsets"`
}

type ZoneChangeStatus struct {
	Name string `json:"name"`
}
