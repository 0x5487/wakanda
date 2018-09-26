package gateway

const (
	PUSH_ROOM = 1 // push to memebers in a room
	PUSH_ALL  = 2 // push to all online members
)

type Job struct {
	OP        int
	RoomID    string
	Command   *Command
	WSMessage *WSMessage
}
