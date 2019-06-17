package command

// ChangeNickname allows to change nickname of a client.
type ChangeNickname struct {
	Nickname string
}

// BlockUser blocks messages from particular client
// identified by the nickname.
type BlockUser struct {
	Nickname string
}

// UnblockUser reverts blockUser action.
type UnblockUser struct {
	Nickname string
}

// JoinRoom command creates chat room on the server and
// connect client to that chat room.
type JoinRoom struct {
	Room string
}

// TextMessage  handles client text input and broadcast
// the message content to all clients connected to the same chat rooms.
type TextMessage struct {
	Message string
}
