package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aquasecurity/trivy/pkg/types"
)

type TrivyScanner struct {
	TimeoutMins int
	ConfigPath  string
}

type TrivyScannerReport struct {
	Results []types.Result
}

func NewTrivyScanner(timeoutMins int, configPath string) *TrivyScanner {
	return &TrivyScanner{
		TimeoutMins: timeoutMins,
		ConfigPath:  configPath,
	}
}

func (ts *TrivyScanner) runCommand(ctx context.Context, cmd string, args ...string) (string, error) {
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

func (cl *TrivyScanner) ParseOutput(dirPath string, resultFile string) ([]*Detail, error) {

	details := make([]*Detail, 0)

	data, err := os.ReadFile(resultFile)
	if err != nil {
		return nil, err
	}

	report := &TrivyScannerReport{}
	err = json.Unmarshal(data, report)
	if err != nil {
		return nil, err
	}

	for _, result := range report.Results {
		fPath := "/" + result.Target
		fName := filepath.Base(fPath)
		info := "secret leak"
		details = append(details, &Detail{
			Name:     fName,
			Path:     fPath,
			Detail:   info,
			Security: UNSAFE_STR,
		})
	}

	return details, nil
}
func (cl *TrivyScanner) ScanDir(ctx context.Context, dirPath string, ignores map[string]bool) (*Result, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cl.TimeoutMins)*time.Minute)
	defer cancel()

	tf, err := os.CreateTemp(os.TempDir(), "trivy-*")
	if err != nil {
		return nil, err
	}
	defer tf.Close()

	args := []string{
		"fs",
		// "--timeout",
		// fmt.Sprintf("%d", cl.TimeoutMins*60),
		"--scanners", "secret",
		"--secret-config", cl.ConfigPath,
		dirPath,
		"-f",
		"json",
		"-o",
		tf.Name(),
	}

	// for k := range ignores {
	// 	args = append(args, "--skip-files strings")
	// 	args = append(args, k)
	// 	args = append(args, "--skip-dirs")
	// 	args = append(args, k)
	// }

	output, err := cl.runCommand(ctx, "/usr/bin/trivy", args...)
	if err != nil {
		fmt.Println(output)
		return nil, err
	}

	details, err := cl.ParseOutput(dirPath, tf.Name())
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
