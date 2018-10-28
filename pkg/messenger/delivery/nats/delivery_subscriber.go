package nats

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/log"
	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"github.com/nats-io/go-nats-streaming"
)

type DeliverySubscriber struct {
	natsConn      stan.Conn
	groupSvc      messenger.GroupServicer
	routerClient  routerProto.RouterServiceClient
	gatewayClient gatewayProto.GatewayServiceClient
}

func NewDeliverySubscriber(natsConn stan.Conn, groupSvc messenger.GroupServicer, routerClient routerProto.RouterServiceClient) *DeliverySubscriber {
	return &DeliverySubscriber{
		natsConn:     natsConn,
		groupSvc:     groupSvc,
		routerClient: routerClient,
	}
}

func (sub *DeliverySubscriber) SubscribeDeliverySubject(ctx context.Context) {
	sub.natsConn.QueueSubscribe("delivery", "worker1", func(m *stan.Msg) {
		log.Debugf("delivery: received a message: %s", string(m.Data))

		ctx := context.Background()
		msg := messenger.Message{}
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			log.Errorf("delivery: json unmarshal message failed: %v", err)
			return
		}

		opts := &messenger.FindGroupMemberOptions{
			GroupID: msg.GroupID,
		}

		groupMembers, err := sub.groupSvc.GroupMembers(ctx, opts)
		if err != nil {
			log.Errorf("delivery: get group members failed: %v", err)
			return
		}

		routeReq := &routerProto.RouteRequest{}
		for _, member := range groupMembers {
			routeReq.MemberIDs = append(routeReq.MemberIDs, member.MemberID)
		}

		routeReply, err := sub.routerClient.Routes(ctx, routeReq)
		if err != nil {
			log.Errorf("delivery: get routes failed: %v", err)
			return
		}

		if len(routeReply.Routes) == 0 {
			return
		}

		// jobReq := gatewayProto.SendJobRequest{}
		// sub.gatewayClient.SendJobs()

		m.Ack()
	}, stan.SetManualAckMode(), stan.DurableName("delivery-remember"))
	log.Info("delivery: delivery subject was subscribed")
}

func (sub *DeliverySubscriber) Shutdown(ctx context.Context) error {
	return nil
}
