package scan

import "context"

type ScanLevel = int

const (
	UNKNOWN ScanLevel = iota + 1
	SUSPICIOUS
	DANGEROUS
)

type ScanResult struct {
	Level  ScanLevel
	Path   string
	Detail interface{}
}

type Scaner interface {
	ScanDir(ctx context.Context, dirPath []string) ([]ScanResult, error)
	ScanFile(ctx context.Context, filePath string) (ScanResult, error)
}
