package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrCacheDown     = errors.New("cache unavailable")
	ErrRepoDown      = errors.New("repository unavailable")
	ErrInvalidSymbol = errors.New("invalid symbol")
	ErrInvalidPeriod = errors.New("invalid period")
)
