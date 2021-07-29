package compilation

import "time"

type Query struct {
	ModulePath string
	BinaryPath string
	Timeout    time.Duration
}
