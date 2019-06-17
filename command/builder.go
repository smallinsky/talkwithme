package command

const (
	nicnameMaxLength     = 32
	userMessageMaxLength = 512
	roomNameMaxLength    = 16
)

var supportedCommands = commandsBuilder{
	"nickname": buildHandler{
		bind: func(args []string) interface{} {
			return ChangeNickname{
				Nickname: args[0],
			}
		},
		validators: []validateFn{
			numberOfArguments(1),
			argumentNMaxLength(0, nicnameMaxLength),
		},
	},
	"block": buildHandler{
		bind: func(args []string) interface{} {
			return BlockUser{
				Nickname: args[0],
			}
		},
		validators: []validateFn{
			numberOfArguments(1),
		},
	},
	"unblock": buildHandler{
		bind: func(args []string) interface{} {
			return UnblockUser{
				Nickname: args[0],
			}
		},
		validators: []validateFn{
			numberOfArguments(1),
		},
	},
	"join": buildHandler{
		bind: func(args []string) interface{} {
			return JoinRoom{
				Room: args[0],
			}
		},
		validators: []validateFn{
			numberOfArguments(1),
		},
	},

	"room": buildHandler{
		bind: func(args []string) interface{} {
			return JoinRoom{
				Room: args[0],
			}
		},
		validators: []validateFn{
			numberOfArguments(1),
			argumentNMaxLength(0, roomNameMaxLength),
		},
	},

	// Wraps client raw text input into command with message argument.
	"textmessage": buildHandler{
		bind: func(args []string) interface{} {
			return TextMessage{
				Message: args[0],
			}
		},
		validators: []validateFn{
			numberOfArguments(1),
			argumentNMaxLength(0, userMessageMaxLength),
		},
	},
}

type buildHandler struct {
	bind       func(args []string) interface{}
	validators []validateFn
}

type validateFn func(args []string) error

type commandsBuilder map[string]buildHandler
