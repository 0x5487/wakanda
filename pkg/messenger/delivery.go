package messenger

import "context"

type DeliveryServicer interface {
	DeliveryMessage(ctx context.Context, msgs []*Message) error
}
