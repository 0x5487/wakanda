package gateway

const (
	OP_PUSH_ROOM = 1 // push to memebers in a room
	OP_PUSH_ALL  = 2 // push to all online members
)

type Job struct {
	OP        int
	RoomID    string
	SessionID uint64
	Command   *Command
	WSMessage *WSMessage
}
