package response

type ZoneDetail struct {
	Account          string        `json:"account"`
	APIRectify       bool          `json:"api_rectify"`
	Catalog          string        `json:"catalog"`
	Dnssec           bool          `json:"dnssec"`
	EditedSerial     int           `json:"edited_serial"`
	ID               string        `json:"id"`
	Kind             string        `json:"kind"`
	LastCheck        int           `json:"last_check"`
	MasterTsigKeyIds []interface{} `json:"master_tsig_key_ids"`
	Masters          []interface{} `json:"masters"`
	Name             string        `json:"name"`
	NotifiedSerial   int           `json:"notified_serial"`
	Nsec3Narrow      bool          `json:"nsec3narrow"`
	Nsec3Param       string        `json:"nsec3param"`
	Rrsets           []struct {
		Comments []interface{} `json:"comments"`
		Name     string        `json:"name"`
		Records  []struct {
			Content  string `json:"content"`
			Disabled bool   `json:"disabled"`
		} `json:"records"`
		TTL  int    `json:"ttl"`
		Type string `json:"type"`
	} `json:"rrsets"`
	Serial          int           `json:"serial"`
	SlaveTsigKeyIds []interface{} `json:"slave_tsig_key_ids"`
	SoaEdit         string        `json:"soa_edit"`
	SoaEditAPI      string        `json:"soa_edit_api"`
	URL             string        `json:"url"`
}

type GetZonesResponse []struct {
	Account        string        `json:"account"`
	Catalog        string        `json:"catalog"`
	Dnssec         bool          `json:"dnssec"`
	EditedSerial   int           `json:"edited_serial"`
	ID             string        `json:"id"`
	Kind           string        `json:"kind"`
	LastCheck      int           `json:"last_check"`
	Masters        []interface{} `json:"masters"`
	Name           string        `json:"name"`
	NotifiedSerial int           `json:"notified_serial"`
	Serial         int           `json:"serial"`
	URL            string        `json:"url"`
}
type GetZonesPdnsAdminResponse struct {
	Account struct {
		Contact     interface{} `json:"contact"`
		Description string      `json:"description"`
		ID          int         `json:"id"`
		Mail        string      `json:"mail"`
		Name        string      `json:"name"`
	} `json:"account"`
	AccountID          int    `json:"account_id"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Status             string `json:"status"`
	UpdateTimeDeactive string `json:"update_time_deactive"`
	UpdatedAt          string `json:"updated_at"`
}

type GetServerStatusResponse struct {
	AutoprimariesURL string `json:"autoprimaries_url"`
	ConfigURL        string `json:"config_url"`
	DaemonType       string `json:"daemon_type"`
	ID               string `json:"id"`
	Type             string `json:"type"`
	URL              string `json:"url"`
	Version          string `json:"version"`
	ZonesURL         string `json:"zones_url"`
}

type ZoneChangeStatus struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}
