package messenger

import "context"

type DeliveryServicer interface {
	DeliveryMessage(ctx context.Context, msg *Message) error
}


