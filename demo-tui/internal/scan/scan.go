package scan

import "context"

type Level = int

const (
	UNKNOWN Level = iota + 1
	SUSPICIOUS
	UNSAFE
)

var (
	UNKNOWN_STR    = "UNKNOWN"
	SUSPICIOUS_STR = "SUSPICIOUS"
	UNSAFE_STR     = "UNSAFE"
)

type Detail struct {
	Name     string
	Security string
	Path     string
	Detail   string
}

type Result struct {
	Count   int
	Stat    map[string]int
	CResult map[string][]Detail
}

func NewEmptyResult() *Result {
	return &Result{
		Count:   0,
		Stat:    map[string]int{},
		CResult: map[string][]Detail{},
	}
}

type Scaner interface {
	Scan(ctx context.Context, dirPath string, ignores map[string]bool) (Result, error)
}

func RunScan(ctx context.Context, dirPath string) ([]Result, error) {

	return nil, nil
}

func MergeResults(a *Result, b *Result) *Result {
	r := NewEmptyResult()
	r.Count = a.Count + b.Count
	r.Stat = a.Stat
	for k := range b.Stat {
		r.Stat[k] += 1
	}
	r.CResult = a.CResult

	for k, details := range b.CResult {
		for _, d := range details {
			if _, ok := r.CResult[k]; !ok {
				r.CResult[k] = make([]Detail, 0)
			}
			r.CResult[k] = append(r.CResult[k], Detail{
				Name:     d.Name,
				Path:     d.Path,
				Security: d.Security,
				Detail:   d.Detail,
			})
		}
	}
	return r
}
