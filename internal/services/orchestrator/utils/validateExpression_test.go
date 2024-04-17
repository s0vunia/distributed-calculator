package orchestratorutils

import "testing"

func Test_validateExpression(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "2+2",
			args: args{expression: "2+2"},
			want: true,
		},
		{
			name: "3*5",
			args: args{expression: "3*5"},
			want: true,
		},
		{
			name: "3/5",
			args: args{expression: "3/5"},
			want: true,
		},
		{
			name: "6-5",
			args: args{expression: "6-5"},
			want: true,
		},
		{
			name: "(6-5)*4",
			args: args{expression: "(6-5)*4"},
			want: true,
		},
		{
			name: "28*9292",
			args: args{expression: "28*9292"},
			want: true,
		},
		{
			name: "10/0",
			args: args{expression: "10/0"},
			want: true,
		},
		{
			name: "10/",
			args: args{expression: "10/"},
			want: false,
		},
		{
			name: "(2+2(",
			args: args{expression: "(2+2("},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateExpression(tt.args.expression); got != tt.want {
				t.Errorf("validateExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
