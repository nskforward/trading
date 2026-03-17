package trading

import (
	"log/slog"

	"github.com/nskforward/trading/types"
)

type Core struct {
	broker types.Broker
	store  *SubscriptionStore
}

func NewCore(broker types.Broker) *Core {
	return &Core{
		broker: broker,
		store:  NewSubscriptionStore(broker),
	}
}

func (c *Core) AddStrategy(strategies ...types.Strategy) error {
	for _, strategy := range strategies {
		err := c.store.AddStrategy(strategy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Core) Run() error {

	err := c.store.Init()
	if err != nil {
		return err
	}

	slog.Info("successfully initialized all strategies")

	return c.store.SubscribeAndWatch()
}
