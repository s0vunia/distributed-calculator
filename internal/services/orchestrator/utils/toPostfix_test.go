package orchestratorutils

import "testing"

func TestInfixToPostfix(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1+2",
			args: args{expression: "1+2"},
			want: "1 2 +",
		},
		{
			name: "(1+2) * 3",
			args: args{expression: "(1+2) * 3"},
			want: "1 2 + 3 *",
		},
		{
			name: "9 + 8 * 2",
			args: args{expression: "9 + 8 * 2"},
			want: "9 8 2 * +",
		},
		{
			name: "1 + 2 * 3 - 4",
			args: args{expression: "1 + 2 * 3 - 4"},
			want: "1 2 3 * + 4 -",
		},
		{
			name: "1 + 5 * 3",
			args: args{expression: "1 + 5 * 3"},
			want: "1 5 3 * +",
		},
		{
			name: "(1 + 5) * 3",
			args: args{expression: "(1 + 5) * 3"},
			want: "1 5 + 3 *",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InfixToPostfix(tt.args.expression); got != tt.want {
				t.Errorf("InfixToPostfix() = %v, want %v", got, tt.want)
			}
		})
	}
}
