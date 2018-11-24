package nats

import (
	"context"
	"encoding/json"

	"github.com/jasonsoft/wakanda/pkg/gateway"

	"github.com/jasonsoft/log"
	gatewayProto "github.com/jasonsoft/wakanda/pkg/gateway/proto"
	"github.com/jasonsoft/wakanda/pkg/messenger"
	routerProto "github.com/jasonsoft/wakanda/pkg/router/proto"
	"github.com/nats-io/go-nats-streaming"
)

type DeliverySubscriber struct {
	natsConn         stan.Conn
	groupSvc         messenger.GroupServicer
	routerClient     routerProto.RouterServiceClient
	gatewayJobClient gatewayProto.JobServiceClient
}

func NewDeliverySubscriber(natsConn stan.Conn, groupSvc messenger.GroupServicer, routerClient routerProto.RouterServiceClient, gatewayJobClient gatewayProto.JobServiceClient) *DeliverySubscriber {
	return &DeliverySubscriber{
		natsConn:         natsConn,
		groupSvc:         groupSvc,
		routerClient:     routerClient,
		gatewayJobClient: gatewayJobClient,
	}
}

func (sub *DeliverySubscriber) SubscribeDeliverySubject(ctx context.Context) {
	sub.natsConn.QueueSubscribe("delivery", "worker1", func(m *stan.Msg) {
		defer func() {
			m.Ack()
		}()

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
			log.Debug("delivery: no routes was found")
			return
		}

		command := &gateway.Command{
			OP:   "MSG",
			Data: m.Data,
		}
		commandBytes, err := json.Marshal(command)
		if err != nil {
			log.Errorf("delivery: json marshal command failed: %v", err)
			return
		}

		jobReq := &gatewayProto.SendJobRequest{}
		for _, route := range routeReply.Routes {
			job := &gatewayProto.Job{
				Type:     "S",
				TargetID: route.SessionID,
				Data:     commandBytes,
			}
			jobReq.Jobs = append(jobReq.Jobs, job)
		}

		_, err = sub.gatewayJobClient.SendJobs(ctx, jobReq)
		if err != nil {
			log.Errorf("delivery: send jobs failed: %v", err)
			return
		}
		log.Debug("delivery: jobs send to gateway")
	}, stan.SetManualAckMode(), stan.DurableName("delivery-remember"))
	log.Info("delivery: delivery subject was subscribed")
}

func (sub *DeliverySubscriber) Shutdown(ctx context.Context) error {
	return nil
}
