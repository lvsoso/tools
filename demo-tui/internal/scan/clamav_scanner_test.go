package scan

import (
	"context"
	"testing"
)

func TestClamAVScanner_ScanDir(t *testing.T) {
	type fields struct {
		TimeoutMins int
		Database    string
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
			name: "test-unsafe",
			fields: fields{
				TimeoutMins: 5,
				Database:    "/home/lv/lvsoso/tools/demo-tui/data/clamavdb",
			},
			args: args{
				ctx:     context.TODO(),
				dirPath: "/home/lv/lvsoso/tools/demo-tui/test/clamav-test/unsafe",
				ignores: map[string]bool{"\\.git": true},
			},
		},
		{
			name: "test-safe",
			fields: fields{
				TimeoutMins: 5,
				Database:    "/home/lv/lvsoso/tools/demo-tui/data/clamavdb",
			},
			args: args{
				ctx:     context.TODO(),
				dirPath: "/home/lv/lvsoso/tools/demo-tui/test/clamav-test/safe",
				ignores: map[string]bool{"\\.git": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &ClamAVScanner{
				TimeoutMins: tt.fields.TimeoutMins,
				Database:    tt.fields.Database,
			}
			got, err := cl.ScanDir(tt.args.ctx, tt.args.dirPath, tt.args.ignores)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClamAVScanner.ScanDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("ClamAVScanner.ScanDir() = %v, want %v", got, tt.want)
			// }
		})
	}
}
