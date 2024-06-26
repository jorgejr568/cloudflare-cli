package cloudflare

import "errors"

var (
	ErrZoneNotFound       = errors.New("zone not found")
	ErrRecordNotFound     = errors.New("record not found")
	ErrRecordDeleteFailed = errors.New("record delete failed")
	ErrRecordAddFailed    = errors.New("record add failed")
	ErrZoneListFailed     = errors.New("zone list failed")
	ErrZoneRecordsFailed  = errors.New("zone records failed")
)
