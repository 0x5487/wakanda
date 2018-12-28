package gateway

import (
	"encoding/json"

	"github.com/jasonsoft/log"
)

type CommandJoinData struct {
	RoomID string `json:"room_id"`
}

func (s *WSSession) handleJoin(commandReq *Command) (*Command, error) {
	log.Debug("gateway: join command is handling")
	joinData := CommandJoinData{}

	err := json.Unmarshal(commandReq.Data, &joinData)
	if err != nil {
		return nil, err
	}

	s.manager.JoinRoom(joinData.RoomID, s)
	return nil, nil
}

type CommandLeaveData struct {
	RoomID string `json:"room_id"`
}

func (s *WSSession) handleLeave(commandReq *Command) (*Command, error) {
	log.Debug("gateway: leave command is handling")
	leaveData := CommandLeaveData{}

	err := json.Unmarshal(commandReq.Data, &leaveData)
	if err != nil {
		return nil, err
	}

	s.manager.LeaveRoom(leaveData.RoomID, s)
	return nil, nil
}

func (s *WSSession) handlePushAll(commandReq *Command) (*Command, error) {
	log.Debugf("gateway: PUSHALL command is handling, session id: %s", s.ID)

	msg := ""
	err := json.Unmarshal(commandReq.Data, &msg)
	if err != nil {
		return nil, err
	}

	// TODO: filter content

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	commandResp := &Command{
		OP:   "MSG",
		Data: data,
	}

	s.manager.PushAll(commandResp)
	return nil, nil
}

func (s *WSSession) handleRoomPush(roomID string, commandReq *Command) (*Command, error) {
	log.Debugf("gateway: PUSHROOM command is handling, session id: %s", s.ID)

	msg := ""
	err := json.Unmarshal(commandReq.Data, &msg)
	if err != nil {
		return nil, err
	}

	// TODO: filter content

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	commandResp := &Command{
		OP:   "MSG",
		Data: data,
	}

	s.manager.PushRoom(roomID, commandResp)
	return nil, nil
}
