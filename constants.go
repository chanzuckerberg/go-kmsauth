package kmsauth

import (
	"time"
)

const (
	timeSkew   = time.Duration(3) * time.Minute
	timeFormat = "%Y%m%dT%H%M%SZ"
)
