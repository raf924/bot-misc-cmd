package pkg

import (
	"github.com/raf924/bot/pkg/domain"
	"reflect"
	"testing"
	"time"
)

func TestMathCommand_Execute(t *testing.T) {
	type args struct {
		argString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Binary Op",
			args:    args{argString: "1 + 3"},
			want:    "4",
			wantErr: false,
		},
		{
			name:    "Multiple Op",
			args:    args{argString: "1 * 2 + 3"},
			want:    "5",
			wantErr: false,
		},
		{
			name:    "Built-in Math",
			args:    args{argString: "Math.sqrt(4)"},
			want:    "2",
			wantErr: false,
		},
		{
			name:    "Floating result",
			args:    args{argString: "7/2"},
			want:    "3.5",
			wantErr: false,
		},
		{
			name:    "Unknown symbol",
			args:    args{argString: "hello"},
			want:    "ReferenceError: hello is not defined at <eval>:1:1(0)",
			wantErr: false,
		},
		{
			name:    "String",
			args:    args{argString: `"hello"`},
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "Unicode",
			args:    args{argString: `"\x20"`},
			want:    " ",
			wantErr: false,
		},
		{
			name:    "Property access",
			args:    args{argString: `""["length"]`},
			want:    "0",
			wantErr: false,
		},
		{
			name:    "This",
			args:    args{argString: "this"},
			want:    "[object global]",
			wantErr: false,
		},
		{
			name:    "Statements",
			args:    args{argString: "1+3;5+4"},
			want:    "9",
			wantErr: false,
		},
		{
			name:    "Loop",
			args:    args{argString: "while(true){}"},
			want:    TimeoutMessage,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MathCommand{}
			got, err := m.Execute(domain.NewCommandMessage(
				"math",
				nil,
				tt.args.argString,
				domain.NewUser("test", "test", domain.RegularUser),
				false,
				time.Now(),
			))
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got[0].Message(), tt.want) {
				t.Errorf("Execute() got = %v, want %v", got[0].Message(), tt.want)
			}
		})
	}
}
