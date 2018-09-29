package gateway

import (
	"encoding/json"

	"github.com/jasonsoft/log"
)

type CommandJoinData struct {
	RoomID string `json:"room_id"`
}

func (s *WSSession) handleJoin(commandReq *Command) (*Command, error) {
	log.Info("gateway: join command is handling")
	joinData := CommandJoinData{}

	err := json.Unmarshal(commandReq.Data, &joinData)
	if err != nil {
		return nil, err
	}

	_manager.JoinRoom(joinData.RoomID, s)
	return nil, nil
}

type CommandLeaveData struct {
	RoomID string `json:"room_id"`
}

func (s *WSSession) handleLeave(commandReq *Command) (*Command, error) {
	log.Info("gateway: leave command is handling")
	leaveData := CommandLeaveData{}

	err := json.Unmarshal(commandReq.Data, &leaveData)
	if err != nil {
		return nil, err
	}

	_manager.LeaveRoom(leaveData.RoomID, s)
	return nil, nil
}

func (s *WSSession) handlePushAll(commandReq *Command) (*Command, error) {
	log.Infof("gateway: push all command is handling, session id: %d", s.ID)

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

	_manager.PushAll(commandResp)
	return nil, nil
}
