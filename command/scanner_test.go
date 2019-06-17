package command

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestCommandScanner(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   interface{}
		err   error
	}{
		{
			name:  "empty input",
			input: "",
			err:   ErrEmptyLine,
		},
		{
			name:  "nickname without arg",
			input: "/nickname",
			err:   ErrInvalidArgsCount,
		},
		{
			name:  "nickname with arg",
			input: "/nickname JohnSnow",
			exp: ChangeNickname{
				Nickname: "JohnSnow",
			},
		},
		{
			name:  "command not supported",
			input: "/notsupported",
			err:   ErrCmdNotSupported,
		},
		{
			name:  "room max length exceeded",
			input: "/room verylongnameoftheroomchannel",
			err:   ErrMaxLengthExeeded,
		},
		{
			name:  "room command",
			input: "/room fun",
			exp: JoinRoom{
				Room: "fun",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Scan(tc.input)
			if errors.Cause(err) != tc.err {
				t.Fatalf("Unexpected error\nGot: %v Expected: %v", err, tc.err)
			}

			if !reflect.DeepEqual(got, tc.exp) {
				t.Fatalf("Commands are not equal\nGot: %+v Expected: %+v", got, tc.exp)
			}
		})
	}
}
