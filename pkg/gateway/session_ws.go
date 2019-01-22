package gateway

import (
	"fmt"
	"context"

	"sync"
	"time"

	"github.com/jasonsoft/wakanda/pkg/identity"

	"github.com/gorilla/websocket"
	"github.com/jasonsoft/log"
	dispatcherProto "github.com/jasonsoft/wakanda/pkg/dispatcher/proto"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"github.com/json-iterator/go"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongPeriod = 60 * time.Second

	// Send pings to peer with this period. Must be less than readWait.
	pingPeriod = 20 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

type WSMessage struct {
	MsgType int
	MsgData []byte
}

type WSSession struct {
	mutex sync.Mutex
	isActive bool
	manager          *Manager
	dispatcherClient dispatcherProto.DispatcherServiceClient
	routerClient     routerProto.RouterServiceClient

	ID          string
	claims      identity.Claims
	socket      *websocket.Conn
	rooms       sync.Map
	roomID      string // member play chatroom and use the roomID
	inChan      chan *WSMessage
	outChan     chan *WSMessage
	commandChan chan *Command
}

func NewWSSession(id string, claims identity.Claims, conn *websocket.Conn, manager *Manager, dispatcherClient dispatcherProto.DispatcherServiceClient, routerClient routerProto.RouterServiceClient, roomID string) *WSSession {
	return &WSSession{
		mutex: sync.Mutex{},
		manager:          manager,
		dispatcherClient: dispatcherClient,
		routerClient:     routerClient,
		isActive: true,
		ID:               id,
		claims:          claims,
		socket:           conn,
		inChan:           make(chan *WSMessage, 1024),
		outChan:          make(chan *WSMessage, 1024),
		commandChan:      make(chan *Command, 1024),
		roomID:           roomID,
	}
}

func (s *WSSession) readLoop() {
	defer func() {
		s.Close()
	}()
	s.socket.SetReadLimit(maxMessageSize)
	//s.socket.SetReadDeadline(time.Now().Add(pongPeriod))
	s.socket.SetPongHandler(func(string) error {
		return nil
	})

	var (
		msgType int
		msgData []byte
		message *WSMessage
		err     error
	)

	for {
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
	//s.socket.SetWriteDeadline(time.Now().Add(writeWait))
	var (
		message *WSMessage
		err     error
	)

	for {
		select {
		case message = <-s.outChan:
			if err = s.socket.WriteMessage(message.MsgType, message.MsgData); err != nil {
				log.Errorf("gateway: wrtieLoop error: %v", err)
				return
			}
			case <-pingTicker.C:
				if err := s.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Errorf("gateway: wrtieLoop ping error: %v", err)
					return
				}
		}
	}
}

func (s *WSSession) commandLoop() {
	var commands []*Command
	timer := time.NewTicker(1 * time.Second).C
	jsoner := jsoniter.ConfigCompatibleWithStandardLibrary

	for {
		select {
		case cmd := <-s.commandChan:
			commands = append(commands, cmd)
		case <-timer:
			//log.Debugf("command chan length: %d", len(s.commandChan))
			if len(commands) > 0 {				
				buf, err := jsoner.Marshal(commands)
				commands = []*Command{}
				if err != nil {
					log.Errorf("gateway: command marshal failed: %v", err)
					continue
				}
				message := &WSMessage{websocket.TextMessage, buf}				
				s.SendMessage(message)				
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
	default:
	}
}

func (s *WSSession) SendCommand(cmd *Command) {
	select {
	case s.commandChan <- cmd:
	default:
	}
}

// Close func which closes websocket session and remove session from bucket and room.
func (s *WSSession) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.isActive {
		s.socket.Close()
		s.manager.DeleteSession(s)
		s.isActive = false
		log.Debugf("gateway: session was closed")
	}
}

func (s *WSSession) refreshRouter() {
	defer func() {
		if r := recover(); r != nil {
			// unknown error
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("gateway: refresh router error: %v", err)
			}
			log.StackTrace().Error(err)
			s.refreshRouter()  // auto-restart
		}
	}()

	timer := time.NewTicker(time.Duration(30) * time.Second)
	log.Debug("gateway: refreshRouter starting")

	for range timer.C {
		in := &routerProto.CreateOrUpdateRouteRequest{
			SessionID:   s.ID,
			MemberID:    s.claims["account_id"].(string),
			GatewayAddr: s.manager.gatewayAddr,
		}
		_, err := s.routerClient.CreateOrUpdateRoute(context.Background(), in)
		if err != nil {
			log.Warnf("gateway: refreshRouter failed: %v", err)
		}
	}
}

func (s *WSSession) StartTasks() {
	defer func() {
		s.Close()
	}()

	s.manager.AddSession(s)

	go s.readLoop()
	go s.writeLoop()
	go s.commandLoop()
	go s.refreshRouter()

	var (
		message     *WSMessage
		commandReq  *Command
		commandResp *Command
		err         error
		//buf         []byte
	)


	for {
		message = s.ReadMessage()

		if message.MsgType != websocket.TextMessage {
			log.Info("gateway: message type is not text message")
			continue
		}

		commandReq, err = CreateCommand(message.MsgData)
		if err != nil {
			log.Errorf("gateway: websocket message is invalid command: %v", err)
			continue
		}
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
				log.Errorf("gateway: handle PUSHALL command error: %v", err)
				continue
			}
		case "PUSHROOM":
			commandResp, err = s.handleRoomPush(s.roomID, commandReq)
			if err != nil {
				log.Errorf("gateway: handle PUSHALL command error: %v", err)
				continue
			}
		default:
			in := &dispatcherProto.DispatcherCommandRequest{
				OP:   commandReq.OP,
				Data: commandReq.Data,
				SenderID: s.claims["account_id"].(string),
				SenderFirstName: s.claims["first_name"].(string),
				SenderLastName: s.claims["last_name"].(string),
			}

			if len(s.roomID) > 0 {
				in.TargetID = s.roomID
			}


			md := metadata.Pairs(
				"req_id", uuid.NewV4().String(),
				"sender_id", s.claims["account_id"].(string),
			)
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			log.Debugf("gateway: handle message: %s", commandReq.OP)
			handleCommandReply, err := s.dispatcherClient.HandleCommand(ctx, in)
			if err != nil {
				log.Errorf("gateway: command error from dispatcher server: %v", err)
				continue
			}

			if handleCommandReply != nil && len(handleCommandReply.OP) > 0 {
				log.Debugf("gateway: receive command resp from server: %s", handleCommandReply.OP)
				commandResp = &Command{
					OP:   handleCommandReply.OP,
					Data: handleCommandReply.Data,
				}
			}
		}

		if commandResp == nil {
			continue
		}

		s.SendCommand(commandResp)

		
	}
}
