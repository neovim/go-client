package nvim

import (
	"reflect"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{
			name:  "Trace",
			level: LogTraceLevel,
			want:  "TraceLevel",
		},
		{
			name:  "Debug",
			level: LogDebugLevel,
			want:  "DebugLevel",
		},
		{
			name:  "Info",
			level: LogInfoLevel,
			want:  "InfoLevel",
		},
		{
			name:  "Warn",
			level: LogWarnLevel,
			want:  "WarnLevel",
		},
		{
			name:  "Error",
			level: LogErrorLevel,
			want:  "ErrorLevel",
		},
		{
			name:  "unknown",
			level: LogLevel(-1),
			want:  "unknown Level",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.level.String(); got != tt.want {
				t.Fatalf("LogLevel.String() = %v, want %v", tt.want, got)
			}
		})
	}
}

func TestUserCommand(t *testing.T) {
	t.Parallel()

	uc := reflect.TypeOf((*UserCommand)(nil)).Elem()

	tests := map[string]struct {
		u UserCommand
	}{
		"UserVimCommand": {
			u: UserVimCommand(""),
		},
		"UserLuaCommand": {
			u: UserLuaCommand{},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := reflect.TypeOf(tt.u)
			if !v.Implements(uc) {
				t.Fatalf("%s type not implements %q interface", v.Name(), "UserCommand")
			}

			tt.u.command()
		})
	}
}
