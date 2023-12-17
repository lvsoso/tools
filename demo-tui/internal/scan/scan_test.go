package scan

import (
	"reflect"
	"testing"
)

func TestMergeResults(t *testing.T) {
	type args struct {
		a *Result
		b *Result
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		// TODO: Add test cases.
		{
			name: "test-merge",
			args: args{
				a: &Result{
					Count: 1,
					Stat: map[string]int{
						UNSAFE_STR: 1,
					},
					CResult: map[string][]Detail{
						UNSAFE_STR: {Detail{Name: "1"}},
					},
				},
				b: &Result{
					Count: 2,
					Stat: map[string]int{
						UNSAFE_STR:     1,
						SUSPICIOUS_STR: 1,
					},
					CResult: map[string][]Detail{
						UNSAFE_STR:     {Detail{Name: "1"}},
						SUSPICIOUS_STR: {Detail{Name: "2"}},
					},
				},
			},
			want: &Result{
				Count: 3,
				Stat: map[string]int{
					UNSAFE_STR:     2,
					SUSPICIOUS_STR: 1,
				},
				CResult: map[string][]Detail{
					UNSAFE_STR:     {Detail{Name: "1"}, Detail{Name: "1"}},
					SUSPICIOUS_STR: {Detail{Name: "2"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeResults(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
