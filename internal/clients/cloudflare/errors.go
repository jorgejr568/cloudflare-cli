package cloudflare

import "errors"

var (
	ErrZoneNotFound        = errors.New("zone not found")
	ErrRecordNotFound      = errors.New("record not found")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrRecordUpdateFailed  = errors.New("record update failed")
	ErrRecordDeleteFailed  = errors.New("record delete failed")
	ErrRecordAddFailed     = errors.New("record add failed")
	ErrZoneListFailed      = errors.New("zone list failed")
	ErrZoneRecordsFailed   = errors.New("zone records failed")
)
