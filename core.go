package trading

import (
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
	return c.store.SubscribeAndWatch()
}
