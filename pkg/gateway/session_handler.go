package gateway

import "encoding/json"

type CommandJoinData struct {
	RoomID string `json:"room_id"`
}

func (s *WSSession) handleJoin(commandReq *Command) (*Command, error) {
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
	leaveData := CommandLeaveData{}

	err := json.Unmarshal(commandReq.Data, &leaveData)
	if err != nil {
		return nil, err
	}

	_manager.LeaveRoom(leaveData.RoomID, s)
	return nil, nil
}
