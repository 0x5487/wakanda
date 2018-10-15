package gateway

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/jasonsoft/wakanda/pkg/dispatcher/proto"

	"github.com/gorilla/websocket"
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/internal/types"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	readWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 5

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type WSSession struct {
	ID      string
	member  *types.Member
	socket  *websocket.Conn
	rooms   sync.Map
	inChan  chan *WSMessage
	outChan chan *WSMessage
}

func NewWSSession(id string, member *types.Member, conn *websocket.Conn) *WSSession {
	return &WSSession{
		ID:      id,
		member:  member,
		socket:  conn,
		inChan:  make(chan *WSMessage, 1024),
		outChan: make(chan *WSMessage, 1024),
	}
}

func (s *WSSession) readLoop() {
	defer func() {
		s.Close()
	}()
	s.socket.SetReadLimit(maxMessageSize)
	s.socket.SetPongHandler(func(string) error { s.socket.SetReadDeadline(time.Now().Add(readWait)); return nil })

	var (
		msgType int
		msgData []byte
		message *WSMessage
		err     error
	)

	for {
		s.socket.SetReadDeadline(time.Now().Add(readWait))
		msgType, msgData, err = s.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				log.Errorf("gateway: websocket message error: %v", err)
			}
			break
		}

		message = &WSMessage{
			MsgType: msgType,
			MsgData: msgData,
		}

		select {
		case s.inChan <- message:
		}
	}
}

func (s *WSSession) writeLoop() {
	defer func() {
		s.Close()
	}()
	pingTicker := time.NewTicker(pingPeriod)

	var (
		message *WSMessage
		err     error
	)
	for {
		select {
		case message = <-s.outChan:
			s.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err = s.socket.WriteMessage(message.MsgType, message.MsgData); err != nil {
				log.Errorf("gateway: wrtieLoop error: %v", err)
				return
			}
		case <-pingTicker.C:
			s.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorf("gateway: wrtieLoop ping error: %v", err)
				return
			}
		}
	}
}

func (s *WSSession) ReadMessage() *WSMessage {
	select {
	case message := <-s.inChan:
		return message
	}
}

func (s *WSSession) SendMessage(msg *WSMessage) {
	select {
	case s.outChan <- msg:
	}
}

func (s *WSSession) Close() {
	s.socket.Close()
	_manager.DeleteSession(s)
	log.Debugf("gateway: session was closed")
}

func (s *WSSession) StartTasks() {
	defer func() {
		s.Close()
	}()

	_manager.AddSession(s)

	go s.readLoop()
	go s.writeLoop()

	var (
		message     *WSMessage
		commandReq  *Command
		commandResp *Command
		err         error
		buf         []byte
	)

	for {
		message = s.ReadMessage()

		if message.MsgType != websocket.TextMessage {
			log.Info("gateway: message type is not text message")
			continue
		}

		commandReq, err = CreateCommand(message.MsgData)
		commandResp = nil

		// handles all commands here
		switch commandReq.OP {
		case "JOIN":
			commandResp, err = s.handleJoin(commandReq)
			if err != nil {
				log.Errorf("gateway: handle JOIN command error: %v", err)
				continue
			}
		case "LEAVE":
			commandResp, err = s.handleLeave(commandReq)
			if err != nil {
				log.Errorf("gateway: handle LEAVE command error: %v", err)
				continue
			}
		case "PUSHALL":
			commandResp, err = s.handlePushAll(commandReq)
			if err != nil {
				log.Errorf("gateway: handle LEAVE command error: %v", err)
				continue
			}
		default:
			in := proto.CommandRequest{
				ReqID: commandReq.RequestID,
				OP:    commandReq.OP,
				Data:  commandReq.Data,
			}

			md := metadata.Pairs("req_id", commandReq.RequestID)
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			handleCommandReply, err := _dispatcherClient.HandleCommand(ctx, &in)
			if err != nil {
				log.Errorf("gateway: command error from messenger server: %v", err)
				continue
			}

			if handleCommandReply != nil && len(handleCommandReply.OP) > 0 {
				log.Debugf("gateway: receive command resp from server: %s", handleCommandReply.OP)
				commandResp = &Command{
					RequestID: handleCommandReply.ReqID,
					OP:        handleCommandReply.OP,
					Data:      handleCommandReply.Data,
				}
			}
		}

		if commandResp != nil {
			buf, err = json.Marshal(*commandResp)
			if err != nil {
				continue
			}

			message = &WSMessage{websocket.TextMessage, buf}
			s.SendMessage(message)
		}
	}
}
