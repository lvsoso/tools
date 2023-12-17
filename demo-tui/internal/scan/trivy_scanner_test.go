package scan

import (
	"context"
	"testing"
)

func TestTrivyScanner_ScanDir(t *testing.T) {
	type fields struct {
		TimeoutMins int
		ConfigPath  string
	}
	type args struct {
		ctx     context.Context
		dirPath string
		ignores map[string]bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Result
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			fields: fields{
				TimeoutMins: 2,
				ConfigPath:  "/home/lv/lvsoso/tools/demo-tui/test/trivy-secret.yaml",
			},
			args: args{
				ctx:     context.Background(),
				dirPath: "/home/lv/lvsoso/tools/demo-tui/test/secret_scan",
				ignores: map[string]bool{".git": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &TrivyScanner{
				TimeoutMins: tt.fields.TimeoutMins,
				ConfigPath:  tt.fields.ConfigPath,
			}
			got, err := cl.ScanDir(tt.args.ctx, tt.args.dirPath, tt.args.ignores)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrivyScanner.ScanDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("TrivyScanner.ScanDir() = %v, want %v", got, tt.want)
			// }
		})
	}
}
