package domain

import (
	"context"
)

type ExchangeClient interface {
	StartTCPClients(ctx context.Context) <-chan PriceTick
}
