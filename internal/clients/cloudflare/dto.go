package cloudflare

type GetZoneByDomainRequest struct {
	Domain string
}

type GetZoneByDomainResponse struct {
	ZoneID string
}

type GetZoneRecordsRequest struct {
	ZoneID string   `json:"-"`
	Name   string   `json:"name,omitempty"`
	Type   ZoneType `json:"type,omitempty"`
}

type ZoneRecord struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       ZoneType               `json:"type"`
	Content    string                 `json:"content"`
	Proxied    bool                   `json:"proxied"`
	Proxiable  bool                   `json:"proxiable"`
	TTL        int                    `json:"ttl"`
	Comment    string                 `json:"comment"`
	ZoneID     string                 `json:"zone_id"`
	ZoneName   string                 `json:"zone_name"`
	Tags       []string               `json:"tags"`
	Locked     bool                   `json:"locked"`
	Meta       map[string]interface{} `json:"meta"`
	CreatedOn  string                 `json:"created_on"`
	ModifiedOn string                 `json:"modified_on"`
}

type GetZoneRecordsResponse struct {
	Records []ZoneRecord `json:"result"`
}

type ZoneRecordRequest struct {
	ID      string   `json:"id,omitempty"`
	Type    ZoneType `json:"type"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Proxied bool     `json:"proxied"`
	TTL     int      `json:"ttl"`
	Tags    []string `json:"tags"`
	Comment string   `json:"comment"`
}

type AddZoneRecordRequest struct {
	ZoneID string
	Record ZoneRecordRequest
}

type AddZoneRecordResponse struct {
	Record ZoneRecord `json:"result"`
}

type DeleteZoneRecordRequest struct {
	ZoneID   string
	RecordID string
}

type UpdateZoneRecordRequest struct {
	ZoneID string
	Record ZoneRecord
}
