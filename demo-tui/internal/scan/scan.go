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
	Sha256 string
}

type Scaner interface {
	Scan(ctx context.Context, dirPath string) ([]ScanResult, error)
}

func RunScan(ctx context.Context, dirPath string) ([]ScanResult, error) {

	return nil, nil
}
