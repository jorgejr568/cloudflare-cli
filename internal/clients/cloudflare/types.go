package cloudflare

import (
	"errors"
	"fmt"
)

type ZoneType string

const (
	ZoneTypeA     ZoneType = "A"
	ZoneTypeAAAA  ZoneType = "AAAA"
	ZoneTypeCNAME ZoneType = "CNAME"
	ZoneTypeTXT   ZoneType = "TXT"
	ZoneTypeSRV   ZoneType = "SRV"
	ZoneTypeMX    ZoneType = "MX"
	ZoneTypeNS    ZoneType = "NS"
	ZoneTypeSOA   ZoneType = "SOA"
)

var (
	InvalidZoneType = errors.New("invalid zone type")
	zoneTypes       = []ZoneType{
		ZoneTypeA,
		ZoneTypeAAAA,
		ZoneTypeCNAME,
		ZoneTypeTXT,
		ZoneTypeSRV,
		ZoneTypeMX,
		ZoneTypeNS,
		ZoneTypeSOA,
	}
)

func ParseZoneType(s string) (ZoneType, error) {
	for _, zt := range zoneTypes {
		if zt == ZoneType(s) {
			return zt, nil
		}
	}
	return "", fmt.Errorf("%w: %s", InvalidZoneType, s)
}
