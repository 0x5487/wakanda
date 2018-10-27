package gateway

const (
	OP_PUSH      = 1
	OP_PUSH_ROOM = 2 // push to memebers in a room
	OP_PUSH_ALL  = 3 // push to all online members
)

type Job struct {
	OP        int
	RoomID    string
	SessionID string
	Command   *Command
	WSMessage *WSMessage
}
