package utils

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestCreateLogger(t *testing.T) {
	type args struct {
		lvl zapcore.Level
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "InfoLevel",
			args: args{
				lvl: zapcore.InfoLevel,
			},
			wantErr: false,
		},
		{
			name: "ErrorLevel",
			args: args{
				lvl: zapcore.ErrorLevel,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateLogger(tt.args.lvl)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("CreateLogger() got = %v, want non-nil value", got)
			}
		})
	}
}
