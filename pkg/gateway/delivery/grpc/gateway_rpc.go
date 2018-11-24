package grpc

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
)

type GatewayServer struct {
	manager *gateway.Manager
}

func NewGatewayServer(manager *gateway.Manager) *GatewayServer {
	return &GatewayServer{
		manager: manager,
	}
}

func (s *GatewayServer) GetServerTime(context.Context, *gatewayProto.Empty) (*gatewayProto.GetServerTimeReply, error) {
	timestamp := time.Now().Unix()
	return &gatewayProto.GetServerTimeReply{
		Time: timestamp,
	}, nil
}

func (s *GatewayServer) CreateCommandStream(stream gatewayProto.GatewayService_CreateCommandStreamServer) error {
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			log.Info("收到客戶端通過context發出的終止信號")
			return ctx.Err()
		default:
			// 接收從客戶端發來的消息
			message, err := stream.Recv()
			if err == io.EOF {
				log.Info("客戶端發送的數據流結束")
				return nil
			}
			if err != nil {
				log.Info("接收數據出錯:", err)
				return err
			}
			// 如果接收正常，則根據接收到的 字符串 執行相應的指令
			switch message.OP {
			case "結束對話\n":
				log.Info("收到'結束對話'指令")
				if err := stream.Send(&gatewayProto.CommandReply{OP: "收到結束指令"}); err != nil {
					return err
				}
				// 收到結束指令時，通過 return nil 終止雙向數據流
				return nil
			case "返回數據流\n":
				log.Info("收到'返回數據流'指令")
				// 收到 收到'返回數據流'指令， 連續返回 10 條數據
				for i := 0; i < 10; i++ {
					if err := stream.Send(&gatewayProto.CommandReply{OP: "數據流 #" + strconv.Itoa(i)}); err != nil {
						return err
					}
				}
			default:
				// 缺省情況下， 返回 '服務端返回: ' + 輸入信息
				log.Infof("[收到消息]: %s", string(message.Data[:]))
				if err := stream.Send(&gatewayProto.CommandReply{OP: "服務端返回: " + message.OP}); err != nil {
					return err
				}
			}
		}
	}
}
