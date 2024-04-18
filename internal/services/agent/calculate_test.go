package agent

import (
	"myproject/internal/config"
	"myproject/internal/models"
	"testing"
)

func TestCalculate(t *testing.T) {
	cfg := config.MustLoadPath("../../config/local_tests.yaml")
	type args struct {
		expression *models.SubExpression
	}
	tests := []struct {
		name    string
		args    args
		wantAns float64
		wantErr bool
	}{
		{
			name: "2+2",
			args: args{
				expression: &models.SubExpression{
					Val1:   2,
					Val2:   2,
					Action: "+",
				},
			},
			wantAns: 4,
			wantErr: false,
		},
		{
			name: "2*2",
			args: args{
				expression: &models.SubExpression{
					Val1:   2,
					Val2:   2,
					Action: "*",
				},
			},
			wantAns: 4,
			wantErr: false,
		},
		{
			name: "100+2",
			args: args{
				expression: &models.SubExpression{
					Val1:   100,
					Val2:   2,
					Action: "+",
				},
			},
			wantAns: 102,
			wantErr: false,
		},
		{
			name: "100/2",
			args: args{
				expression: &models.SubExpression{
					Val1:   100,
					Val2:   2,
					Action: "/",
				},
			},
			wantAns: 50,
			wantErr: false,
		},
		{
			name: "1000-50",
			args: args{
				expression: &models.SubExpression{
					Val1:   1000,
					Val2:   50,
					Action: "-",
				},
			},
			wantAns: 950,
			wantErr: false,
		},
		{
			name: "1000/0",
			args: args{
				expression: &models.SubExpression{
					Val1:   1000,
					Val2:   0,
					Action: "/",
				},
			},
			wantAns: 0,
			wantErr: true,
		},
		{
			name: "&*100/",
			args: args{
				expression: &models.SubExpression{
					Val1:   100,
					Val2:   2,
					Action: "*&*",
				},
			},
			wantAns: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAns, err := Calculate(tt.args.expression, cfg.CalculationTimeouts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAns != tt.wantAns {
				t.Errorf("Calculate() gotAns = %v, want %v", gotAns, tt.wantAns)
			}
		})
	}
}
