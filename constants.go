package kmsauth

import (
	"time"
)

const (
	//TimeFormat the time format for kmsauth tokens
	TimeFormat = time.RFC3339
	// timeSkew how much to compensate for time skew
	timeSkew = time.Duration(3) * time.Minute
)
