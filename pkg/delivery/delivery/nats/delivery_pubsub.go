package nats

import (
	"context"

	"github.com/gogo/protobuf/proto"

	"github.com/jasonsoft/log"
	deliveryProto "github.com/jasonsoft/wakanda/pkg/delivery/proto"
	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"github.com/nats-io/go-nats-streaming"
)

type DeliveryPubSub struct {
	natsConn         stan.Conn
	groupSvc         messenger.GroupServicer
	routerClient     routerProto.RouterServiceClient
	gatewayJobClient gatewayProto.JobServiceClient
}

func NewDeliveryPubSub(natsConn stan.Conn, groupSvc messenger.GroupServicer, routerClient routerProto.RouterServiceClient, gatewayJobClient gatewayProto.JobServiceClient) *DeliveryPubSub {
	return &DeliveryPubSub{
		natsConn:         natsConn,
		groupSvc:         groupSvc,
		routerClient:     routerClient,
		gatewayJobClient: gatewayJobClient,
	}
}

func (ps *DeliveryPubSub) SubscribeMSGChatroom(ctx context.Context) {
	ps.natsConn.QueueSubscribe("delivery-chatroom", "worker1", func(m *stan.Msg) {
		defer func() {
			m.Ack()
		}()

		//log.Debugf("delivery: received a message: %s", string(m.Data))

		cmdDelivery := deliveryProto.DeliveryChatroomMessageCommand{}
		err := proto.Unmarshal(m.Data, &cmdDelivery)
		if err != nil {
			log.Errorf("delivery: proto unmarshal message failed: %v", err)
			return
		}

		ctx := context.Background()
		// send Job requests to gateway and ask gateway to send message to member
		jobReq := &gatewayProto.SendJobRequest{}
		job := &gatewayProto.Job{
			Type:     "PUSHROOM",
			TargetID: cmdDelivery.RoomID,
			Data:     cmdDelivery.Data,
		}
		jobReq.Jobs = append(jobReq.Jobs, job)

		_, err = ps.gatewayJobClient.SendJobs(ctx, jobReq)
		if err != nil {
			log.Errorf("delivery: send jobs failed: %v", err)
			return
		}
		log.Debug("delivery: jobs send to gateway")
	}, stan.SetManualAckMode(), stan.DurableName("msg-chatroom-remember"))
	log.Info("delivery: msg-chatroom subject was subscribed")
}

func (ps *DeliveryPubSub) Shutdown(ctx context.Context) error {
	return nil
}
