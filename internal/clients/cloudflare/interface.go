package cloudflare

import "context"

type CloudflareClient interface {
	GetZoneByDomain(context.Context, GetZoneByDomainRequest) (*GetZoneByDomainResponse, error)
	GetZoneRecords(context.Context, GetZoneRecordsRequest) (*GetZoneRecordsResponse, error)
	AddZoneRecord(context.Context, AddZoneRecordRequest) (*AddZoneRecordResponse, error)
	DeleteZoneRecord(context.Context, DeleteZoneRecordRequest) error
}
