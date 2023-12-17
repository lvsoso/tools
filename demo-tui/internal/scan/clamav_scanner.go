package scan

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ClamAVScanner struct {
	TimeoutMins int
	Database    string
}

func NewClamAVScanner(timeoutMins int, database int) *ClamAVScanner {
	return &ClamAVScanner{
		TimeoutMins: timeoutMins,
		Database:    "/home/lv/lvsoso/tools/demo-tui/data/clamavdb",
	}
}

func (cl *ClamAVScanner) runCommand(ctx context.Context, cmd string, args ...string) (string, error) {
	c := exec.CommandContext(ctx, cmd, args...)
	output, err := c.CombinedOutput()
	if err != nil && c.Err != nil {
		return string(output), err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command %s timed out", cmd)
	}

	if ctx.Err() == context.Canceled {
		return "", fmt.Errorf("command %s canceled", cmd)
	}
	return string(output), nil
}

func (cl *ClamAVScanner) ParseOutput(dirPath string, output string) ([]*Detail, error) {
	lines := strings.Split(output, "\n")
	details := make([]*Detail, 0)

	for _, line := range lines {
		if strings.Contains(line, "SCAN SUMMARY") {
			return details, nil
		}

		pathAndResult := strings.Split(line, ":")
		if len(pathAndResult) < 2 {
			continue
		}

		if strings.Contains(pathAndResult[1], "OK") {
		} else {
			fPath := strings.Replace(pathAndResult[0], dirPath, "", 1)
			fName := filepath.Base(fPath)
			info := strings.TrimSpace(strings.TrimRight(pathAndResult[1], "FOUND"))
			details = append(details, &Detail{
				Name:     fName,
				Path:     fPath,
				Detail:   info,
				Security: UNSAFE_STR,
			})
		}
	}

	return details, nil
}
func (cl *ClamAVScanner) ScanDir(ctx context.Context, dirPath string, ignores map[string]bool) (*Result, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cl.TimeoutMins)*time.Minute)
	defer cancel()

	// tmp file
	// tf, err := os.CreateTemp(os.TempDir(), "clamavscan-*")
	// if err != nil {
	// 	return nil, err
	// }
	// defer tf.Close()

	args := []string{
		"-d",
		cl.Database,
		"-r",
		dirPath,
		// "--stdout",
		// tf.Name(),
	}

	for k := range ignores {
		args = append(args, "--exclude")
		args = append(args, k)
		args = append(args, "--exclude-dir")
		args = append(args, k)
	}

	output, err := cl.runCommand(ctx, "/usr/bin/clamscan", args...)
	if err != nil {
		return nil, err
	}

	details, err := cl.ParseOutput(dirPath, output)
	if err != nil {
		return nil, err
	}

	result := NewEmptyResult()
	for _, d := range details {
		result.Count += 1
		if _, ok := result.Stat[d.Security]; !ok {
			result.Stat[d.Security] = 0
		}
		result.Stat[d.Security] += 1
		if _, ok := result.CResult[d.Security]; !ok {
			result.CResult[d.Security] = []Detail{}
		}
		result.CResult[d.Security] = append(result.CResult[d.Security],
			Detail{
				Name:     d.Name,
				Path:     d.Path,
				Security: d.Security,
				Detail:   d.Detail,
			})
	}
	return result, nil
}
